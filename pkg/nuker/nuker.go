package nuker

import (
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
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

	go func() {
		defer close(results)

		for {
			deletedCount := 0

			for _, service := range services {
				if service.IsGlobal() {
					deleted, err := deleteResourcesForServiceInRegion(service, regions[0], currentAccount)
					if err != nil {
						results <- NukeResult{Success: false, Error: err}
					} else {
						deletedCount += len(deleted)
						for _, resource := range deleted {
							results <- NukeResult{Success: true, Resource: resource}
						}
					}
				} else {
					for _, region := range regions {
						deleted, err := deleteResourcesForServiceInRegion(service, region, currentAccount)
						if err != nil {
							results <- NukeResult{Success: false, Error: err}
						} else {
							for _, resource := range deleted {
								results <- NukeResult{Success: true, Resource: resource}
							}
							deletedCount += len(deleted)
						}
					}
				}
			}

			if deletedCount == 0 {
				break
			}
		}
	}()

	return results
}

func deleteResourcesForServiceInRegion(service cloud.Service, region account.Region, currentAccount account.Account) ([]cloud.Resource, error) {
	deletedResources := make([]cloud.Resource, 0)
	foundResources, err := service.List(region, currentAccount)
	if err != nil {
		return nil, err
	} else {
		for _, resource := range foundResources {
			err := resource.Delete(region, currentAccount)
			if err != nil {
				return nil, err
			} else {
				deletedResources = append(deletedResources, resource)
			}
		}
	}
	return deletedResources, nil
}
