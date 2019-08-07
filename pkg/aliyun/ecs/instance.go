package ecs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type Instances struct{}

type Instance struct {
	ecs.Instance
}

// String outputs name of the service
func (s Instances) String() string {
	return "ECS instance"
}

// List returns a list of all machines
func (s Instances) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
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
		instances = append(instances, Instance{Instance: instance})
	}

	return instances, nil
}

func (e Instance) String() string {
	return e.InstanceId
}

func (e Instance) Delete(region account.Region, account account.Account) error {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := ecs.CreateDeleteInstanceRequest()
	request.InstanceId = e.InstanceId
	request.Force = "true"
	request.TerminateSubscription = "true"

	_, err = client.DeleteInstance(request)
	if err != nil {
		return err
	}

	return nil
}