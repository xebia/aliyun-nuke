package nuker

import (
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/aliyun"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

// Nuke actually removes all resources in a loop. It will keep on going until no resources
// were deleted any more.
func Nuke(account account.Account) []cloud.Resource {
	services := []cloud.Service{
		aliyun.OssService{},
	}

	deletedResources := make([]cloud.Resource, 0)
	deleted := 1
	for deleted > 0 {
		deleted = 0

		for _, service := range services {
			foundResources, _ := service.List(account)
			for _, resource := range foundResources {
				ok, _ := resource.Delete()
				if ok {
					deletedResources = append(deletedResources, resource)
					deleted++
				}
			}
		}
	}
	return deletedResources
}
