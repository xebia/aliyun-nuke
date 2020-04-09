package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type EcsSecurityGroups struct{}

type EcsSecurityGroup struct {
	ecs.SecurityGroup
}

func init() {
	cloud.RegisterService(EcsSecurityGroups{})
}

func (s EcsSecurityGroups) IsGlobal() bool {
	return false
}

func (s EcsSecurityGroups) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := ecs.CreateDescribeSecurityGroupsRequest()
	request.PageSize = "50"
	response, err := client.DescribeSecurityGroups(request)
	if err != nil {
		return nil, err
	}

	groups := make([]cloud.Resource, 0)
	for _, securityGroup := range response.SecurityGroups.SecurityGroup {
		groups = append(groups, EcsSecurityGroup{SecurityGroup: securityGroup})
	}

	return groups, nil
}

func (s EcsSecurityGroup) Id() string {
	return s.SecurityGroupId
}

func (s EcsSecurityGroup) Type() string {
	return "Security group"
}

func (s EcsSecurityGroup) Delete(region account.Region, account account.Account) error {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := ecs.CreateDeleteSecurityGroupRequest()
	request.SecurityGroupId = s.SecurityGroupId
	_, err = client.DeleteSecurityGroup(request)
	if err != nil {
		return err
	}

	return nil
}
