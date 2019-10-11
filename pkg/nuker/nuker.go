package nuker

import (
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/aliyun/ecs"
	"github.com/xebia/aliyun-nuke/pkg/aliyun/oss"
	"github.com/xebia/aliyun-nuke/pkg/aliyun/ram"
	"github.com/xebia/aliyun-nuke/pkg/aliyun/vpc"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

// NukeItAll will nuke (delete) all Alibaba Cloud services in all regions
func NukeItAll(currentAccount account.Account) ([]cloud.Resource, error) {
	services := []cloud.Service{
		oss.Buckets{},
		ecs.Instances{},
		ecs.SecurityGroups{},
		ram.Groups{},
		ram.Policies{},
		ram.Roles{},
		ram.Users{},
		vpc.Vpcs{},
		vpc.VSwitches{},
	}
	return Nuke(currentAccount, services, account.Regions)
}

// Nuke removes all resources of specified services in specified regions in a loop.
// It will keep on going until no resources were deleted any more.
func Nuke(currentAccount account.Account, services []cloud.Service, regions []account.Region) ([]cloud.Resource, error) {
	deletedResources := make([]cloud.Resource, 0)

	for {
		deletedCount := 0

		for _, service := range services {
			if service.IsGlobal() {
				deleted, _ := deleteResourcesForServiceInRegion(service, regions[0], currentAccount)
				deletedResources = append(deletedResources, deleted...)
				deletedCount += len(deleted)
			} else {
				for _, region := range regions {
					deleted, _ := deleteResourcesForServiceInRegion(service, region, currentAccount)
					deletedResources = append(deletedResources, deleted...)
					deletedCount += len(deleted)
				}
			}
		}

		if deletedCount == 0 {
			break
		}
	}

	return deletedResources, nil
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
