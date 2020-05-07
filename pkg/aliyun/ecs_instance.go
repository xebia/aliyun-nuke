package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type EcsInstances struct{}

type EcsInstance struct {
	ecs.Instance
}

func init() {
	cloud.RegisterService(EcsInstances{})
}

func (i EcsInstances) IsGlobal() bool {
	return false
}

func (i EcsInstances) List(region account.Region, account account.Account, force bool) ([]cloud.Resource, error) {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := ecs.CreateDescribeInstancesRequest()
	request.PageSize = "99"
	response, err := client.DescribeInstances(request)
	if err != nil {
		return nil, err
	}

	instances := make([]cloud.Resource, 0)
	for _, instance := range response.Instances.Instance {
		instances = append(instances, EcsInstance{Instance: instance})
	}

	return instances, nil
}

func (i EcsInstance) Id() string {
	return i.InstanceId
}

func (i EcsInstance) Type() string {
	return "ECS instance"
}

func (i EcsInstance) Delete(region account.Region, account account.Account) error {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := ecs.CreateDeleteInstanceRequest()
	request.InstanceId = i.InstanceId
	request.Force = "true"
	request.TerminateSubscription = "true"

	_, err = client.DeleteInstance(request)
	if err != nil {
		return err
	}

	return nil
}
