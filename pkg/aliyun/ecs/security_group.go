package ecs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/xebia/aliyun-nuke/pkg/nuker"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type SecurityGroups struct{}

type SecurityGroup struct {
	ecs.SecurityGroup
}

func init() {
	cloud.RegisterService(SecurityGroups{})
}

func (s SecurityGroups) IsGlobal() bool {
	return false
}

func (s SecurityGroups) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
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
		groups = append(groups, SecurityGroup{SecurityGroup: securityGroup})
	}

	return groups, nil
}

func (s SecurityGroup) Id() string {
	return s.SecurityGroupId
}

func (s SecurityGroup) Type() string {
	return "Security group"
}

func (s SecurityGroup) Delete(region account.Region, account account.Account) error {
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
