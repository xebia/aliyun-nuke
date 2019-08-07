package nuker

import (
	"fmt"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/aliyun/ecs"
	"github.com/xebia/aliyun-nuke/pkg/aliyun/oss"
	"github.com/xebia/aliyun-nuke/pkg/aliyun/vpc"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

// Nuke actually removes all resources in a loop. It will keep on going until no resources
// were deleted any more.
func Nuke(currentAccount account.Account) []cloud.Resource {
	services := []cloud.Service{
		oss.Buckets{},
		ecs.Instances{},
		ecs.SecurityGroups{},
		vpc.Vpcs{},
		vpc.VSwitches{},
	}

	deletedResources := make([]cloud.Resource, 0)
	deleted := 1
	for deleted > 0 {
		deleted = 0

		for _, service := range services {
			for _, region := range account.Regions {
				foundResources, err := service.List(region, currentAccount)
				if err != nil {
					fmt.Println(err)
				} else {
					for _, resource := range foundResources {
						err := resource.Delete(region, currentAccount)
						if err != nil {
							fmt.Println(err)
						} else {
							deletedResources = append(deletedResources, resource)
							deleted++
						}
					}
				}
			}
		}
	}
	return deletedResources
}
