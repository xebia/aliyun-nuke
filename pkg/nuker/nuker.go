package nuker

import (
	"github.com/xebia/aliyun-nuke/pkg/account"
	_ "github.com/xebia/aliyun-nuke/pkg/aliyun"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
	"reflect"
	"time"
)

type NukeResult struct {
	Resource cloud.Resource
	Success  bool
	Skipped  bool
	Error    error
}

// NukeItAll will nuke (delete) all Alibaba Cloud services in the specified regions
func NukeItAll(currentAccount account.Account, regions []account.Region, excludedIds []string, force bool) <-chan NukeResult {
	return Nuke(currentAccount, cloud.Services, regions, excludedIds, force)
}

// Nuke removes all resources of specified services in specified regions in a loop.
// It will keep on going until no resources were deleted any more.
func Nuke(currentAccount account.Account, services []cloud.Service, regions []account.Region, excludedIds []string, force bool) <-chan NukeResult {
	results := make(chan NukeResult)

	emptyServices := make([]string, len(services))
	emptyRegionsPerService := make(map[string][]string)

	go func() {
		defer close(results)

		maxRetries := 60
		currentRetry := 0
		for {
			totalLeftOverCount := 0

			for _, service := range services {
				serviceType := reflect.TypeOf(service).String()
				if !elementIn(emptyServices, serviceType) {
					serviceLeftOverCount := 0

					if service.IsGlobal() {
						found, deleted, skipped, errors := deleteResourcesForServiceInRegion(service, "eu-central-1", currentAccount, excludedIds, force)

						leftOvers := len(found) - len(deleted) - len(skipped)
						serviceLeftOverCount += leftOvers
						totalLeftOverCount += leftOvers

						for _, resource := range deleted {
							results <- NukeResult{Success: true, Resource: resource}
						}

						for _, err := range errors {
							results <- NukeResult{Success: false, Error: err}
						}

						for _, skipped := range skipped {
							results <- NukeResult{Skipped: true, Resource: skipped}
						}
					} else {
						for _, region := range regions {
							if !elementIn(emptyRegionsPerService[serviceType], string(region)) {
								found, deleted, skipped, errors := deleteResourcesForServiceInRegion(service, region, currentAccount, excludedIds, force)

								leftOvers := len(found) - len(deleted) - len(skipped)
								serviceLeftOverCount += leftOvers
								totalLeftOverCount += leftOvers

								for _, resource := range deleted {
									results <- NukeResult{Success: true, Resource: resource}
								}

								for _, err := range errors {
									results <- NukeResult{Success: false, Error: err}
								}

								for _, skipped := range skipped {
									results <- NukeResult{Skipped: true, Resource: skipped}
								}

								if (leftOvers) < 1 {
									// Remove this region for this service, as the region is empty
									emptyRegionsPerService[serviceType] = append(emptyRegionsPerService[serviceType], string(region))
								}
							}
						}
					}

					if serviceLeftOverCount < 1 {
						// Remove this service completely, as no resources in any region
						emptyServices = append(emptyServices, serviceType)
					}
				}
			}

			if totalLeftOverCount == 0 || currentRetry == maxRetries {
				break
			} else {
				// Sleep to allow some time for deletion
				time.Sleep(1 * time.Second)

				currentRetry++
			}
		}
	}()

	return results
}

func elementIn(elements []string, element string) bool {
	for _, item := range elements {
		if item == element {
			return true
		}
	}
	return false
}

func deleteResourcesForServiceInRegion(service cloud.Service, region account.Region, currentAccount account.Account, excludedIds []string, force bool) ([]cloud.Resource, []cloud.Resource, []cloud.Resource, []error) {
	foundResources, err := service.List(region, currentAccount, force)

	if err != nil {
		return nil, nil, nil, []error{err}
	}

	deletedResources := make([]cloud.Resource, 0)
	skippedResources := make([]cloud.Resource, 0)
	errors := make([]error, 0)
	for _, resource := range foundResources {
		if !elementIn(excludedIds, resource.Id()) {
			err := resource.Delete(region, currentAccount)
			if err != nil {
				errors = append(errors, err)
			} else {
				deletedResources = append(deletedResources, resource)
			}
		} else {
			skippedResources = append(skippedResources, resource)
		}
	}

	return foundResources, deletedResources, skippedResources, errors
}
