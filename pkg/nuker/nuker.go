package nuker

import (
	"github.com/xebia/aliyun-nuke/pkg/account"
	_ "github.com/xebia/aliyun-nuke/pkg/aliyun"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
	"reflect"
	"time"
)

type NukeResult struct {
	Success  bool
	Resource cloud.Resource
	Error    error
}

// NukeItAll will nuke (delete) all Alibaba Cloud services in the specified regions
func NukeItAll(currentAccount account.Account, regions []account.Region) <-chan NukeResult {
	return Nuke(currentAccount, cloud.Services, regions)
}

// Nuke removes all resources of specified services in specified regions in a loop.
// It will keep on going until no resources were deleted any more.
func Nuke(currentAccount account.Account, services []cloud.Service, regions []account.Region) <-chan NukeResult {
	results := make(chan NukeResult)

	emptyServices := make([]string, len(services))
	emptyRegionsPerService := make(map[string][]string)

	go func() {
		defer close(results)

		maxRetries := 10
		currentRetry := 0
		for {
			totalLeftOverCount := 0

			for _, service := range services {
				serviceType := reflect.TypeOf(service).String()
				if !elementIn(emptyServices, serviceType) {
					serviceLeftOverCount := 0

					if service.IsGlobal() {
						found, deleted, err := deleteResourcesForServiceInRegion(service, "eu-central-1", currentAccount)
						leftOvers := len(found) - len(deleted)
						serviceLeftOverCount += leftOvers

						if err != nil {
							results <- NukeResult{Success: false, Error: err}
						} else {
							totalLeftOverCount += leftOvers
							for _, resource := range deleted {
								results <- NukeResult{Success: true, Resource: resource}
							}
						}
					} else {
						for _, region := range regions {
							if !elementIn(emptyRegionsPerService[serviceType], string(region)) {
								found, deleted, err := deleteResourcesForServiceInRegion(service, region, currentAccount)

								leftOvers := len(found) - len(deleted)
								serviceLeftOverCount += leftOvers
								totalLeftOverCount += leftOvers

								if err != nil {
									results <- NukeResult{Success: false, Error: err}
								} else {
									for _, resource := range deleted {
										results <- NukeResult{Success: true, Resource: resource}
									}
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

func deleteResourcesForServiceInRegion(service cloud.Service, region account.Region, currentAccount account.Account) ([]cloud.Resource, []cloud.Resource, error) {
	foundResources, err := service.List(region, currentAccount)
	if err != nil {
		return nil, nil, err
	}

	deletedResources := make([]cloud.Resource, 0)
	for _, resource := range foundResources {
		err := resource.Delete(region, currentAccount)
		if err != nil {
			return foundResources, deletedResources, err
		} else {
			deletedResources = append(deletedResources, resource)
		}
	}

	return foundResources, deletedResources, nil
}
