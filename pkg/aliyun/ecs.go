package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

//Elastic Compute Service

type EcsService struct{}

type EcsResource struct {
	InstanceId string

	Region  account.Region
	Account account.Account
}

// String outputs name of the service
func (s EcsService) String() string {
	return "ECS"
}

// List returns a list of all machines
func (s EcsService) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	instanceRequest := ecs.CreateDescribeInstancesRequest()
	instanceRequest.PageSize = "99"
	instancesResponse, err := client.DescribeInstances(instanceRequest)
	if err != nil {
		return nil, err
	}

	instances := make([]cloud.Resource, 0)
	for _, instance := range instancesResponse.Instances.Instance {
		instances = append(instances, EcsResource{InstanceId: instance.InstanceId, Region: region, Account: account})
	}

	return instances, nil
}

func (e EcsResource) String() string {
	return e.InstanceId
}

func (e EcsResource) Delete() error {
	client, err := ecs.NewClientWithAccessKey(string(e.Region), e.Account.AccessKeyID, e.Account.AccessKeySecret)
	if err != nil {
		return err
	}

	deleteInstanceRequest := ecs.CreateDeleteInstanceRequest()
	deleteInstanceRequest.InstanceId = e.InstanceId
	deleteInstanceRequest.Force = "true"
	deleteInstanceRequest.TerminateSubscription = "true"

	_, err = client.DeleteInstance(deleteInstanceRequest)
	if err != nil {
		return err
	}

	return nil
}
