package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type RdsInstances struct{}

type RdsInstance struct {
	rds.DBInstanceInDescribeDBInstances
}

func init() {
	cloud.RegisterService(RdsInstances{})
}

func (r RdsInstances) IsGlobal() bool {
	return false
}

func (r RdsInstances) List(region account.Region, account account.Account, force bool) ([]cloud.Resource, error) {
	client, err := rds.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := rds.CreateDescribeDBInstancesRequest()
	request.PageSize = "100"
	response, err := client.DescribeDBInstances(request)
	if err != nil {
		return nil, err
	}

	instances := make([]cloud.Resource, 0)
	for _, instance := range response.Items.DBInstance {
		instances = append(instances, RdsInstance{DBInstanceInDescribeDBInstances: instance})
	}

	return instances, nil
}

func (r RdsInstance) Id() string {
	return r.DBInstanceId
}

func (r RdsInstance) Type() string {
	return "RDS instance"
}

func (r RdsInstance) Delete(region account.Region, account account.Account) error {
	client, err := rds.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := rds.CreateDeleteDBInstanceRequest()
	request.DBInstanceId = r.DBInstanceId

	_, err = client.DeleteDBInstance(request)
	if err != nil {
		return err
	}

	return nil
}
