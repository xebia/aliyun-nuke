package vpc

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type Eips struct{}

type Eip struct {
	vpc.EipAddress
}

func init() {
	cloud.RegisterService(Eips{})
}

func (e Eips) IsGlobal() bool {
	return false
}

func (e Eips) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
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
		eips = append(eips, Eip{EipAddress:eip})
	}
	return eips, nil
}

func (e Eip) Id() string {
	return e.AllocationId
}

func (e Eip) Type() string {
	return "EIP"
}

func (e Eip) Delete(region account.Region, account account.Account) error {
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
