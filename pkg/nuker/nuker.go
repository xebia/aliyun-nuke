package nuker

import (
	"fmt"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/aliyun"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

// Nuke actually removes all resources in a loop. It will keep on going until no resources
// were deleted any more.
func Nuke(currentAccount account.Account) []cloud.Resource {
	services := []cloud.Service{
		//aliyun.OssService{},
		aliyun.EcsService{},
	}

	deletedResources := make([]cloud.Resource, 0)
	deleted := 1
	for deleted > 0 {
		deleted = 0

		for _, service := range services {
			for _, region := range account.Regions {
				foundResources, _ := service.List(region, currentAccount)
				for _, resource := range foundResources {
					err := resource.Delete()
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
	return deletedResources
}
