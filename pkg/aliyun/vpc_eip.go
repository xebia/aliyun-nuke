package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type VpcEips struct{}

type VpcEip struct {
	vpc.EipAddress
}

func init() {
	cloud.RegisterService(VpcEips{})
}

func (e VpcEips) IsGlobal() bool {
	return false
}

func (e VpcEips) List(region account.Region, account account.Account, force bool) ([]cloud.Resource, error) {
	client, err := vpc.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := vpc.CreateDescribeEipAddressesRequest()
	request.PageSize = "50"
	response, err := client.DescribeEipAddresses(request)
	if err != nil {
		return nil, err
	}

	eips := make([]cloud.Resource, 0)
	for _, eip := range response.EipAddresses.EipAddress {
		eips = append(eips, VpcEip{EipAddress: eip})
	}
	return eips, nil
}

func (e VpcEip) Id() string {
	return e.AllocationId
}

func (e VpcEip) Type() string {
	return "EIP"
}

func (e VpcEip) Delete(region account.Region, account account.Account) error {
	client, err := vpc.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := vpc.CreateReleaseEipAddressRequest()
	request.AllocationId = e.AllocationId

	_, err = client.ReleaseEipAddress(request)
	if err != nil {
		return err
	}

	return nil
}
