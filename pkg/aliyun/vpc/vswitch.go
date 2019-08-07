package vpc

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type VSwitches struct {}

type VSwitch struct {
	vpc.VSwitch
}

// String outputs name of the service
func (v VSwitches) String() string {
	return "ECS instance"
}

// List returns a list of all machines
func (v VSwitches) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := vpc.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := vpc.CreateDescribeVSwitchesRequest()
	request.PageSize = "50"
	response, err := client.DescribeVSwitches(request)
	if err != nil {
		return nil, err
	}

	vswitches := make([]cloud.Resource, 0)
	for _, vswitch := range response.VSwitches.VSwitch {
		vswitches = append(vswitches, VSwitch{VSwitch: vswitch})
	}

	return vswitches, nil
}

func (v VSwitch) String() string {
	return v.VSwitchId
}

func (v VSwitch) Delete(region account.Region, account account.Account) error {
	client, err := vpc.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := vpc.CreateDeleteVSwitchRequest()
	request.VSwitchId = v.VSwitchId

	_, err = client.DeleteVSwitch(request)
	if err != nil {
		return err
	}

	return nil
}