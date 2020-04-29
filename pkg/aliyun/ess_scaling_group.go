package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type EssScalingGroups struct{}

type EssScalingGroup struct {
	ess.ScalingGroup
}

func init() {
	cloud.RegisterService(EssScalingGroups{})
}

func (s EssScalingGroups) IsGlobal() bool {
	return false
}

func (s EssScalingGroups) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := ess.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := ess.CreateDescribeScalingGroupsRequest()
	request.PageSize = "50"
	response, err := client.DescribeScalingGroups(request)
	if err != nil {
		return nil, err
	}

	scalingGroups := make([]cloud.Resource, 0)
	for _, scalingGroup := range response.ScalingGroups.ScalingGroup {
		scalingGroups = append(scalingGroups, EssScalingGroup{ScalingGroup: scalingGroup})
	}

	return scalingGroups, nil
}

func (s EssScalingGroup) Id() string {
	return s.ScalingGroupId
}

func (s EssScalingGroup) Type() string {
	return "ESS scaling group"
}

func (s EssScalingGroup) Delete(region account.Region, account account.Account) error {
	client, err := ess.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := ess.CreateDeleteScalingGroupRequest()
	request.ScalingGroupId = s.ScalingGroupId
	request.ForceDelete = "true"

	_, err = client.DeleteScalingGroup(request)
	if err != nil {
		return err
	}

	return nil
}
