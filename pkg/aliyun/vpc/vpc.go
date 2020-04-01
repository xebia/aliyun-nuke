package vpc

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type Vpcs struct{}

type Vpc struct {
	vpc.Vpc
}

func (v Vpcs) IsGlobal() bool {
	return false
}

func (v Vpcs) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := vpc.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := vpc.CreateDescribeVpcsRequest()
	request.PageSize = "50"
	response, err := client.DescribeVpcs(request)
	if err != nil {
		return nil, err
	}

	vpcs := make([]cloud.Resource, 0)
	for _, vpcItem := range response.Vpcs.Vpc {
		vpcs = append(vpcs, Vpc{Vpc: vpcItem})
	}

	return vpcs, nil
}

func (v Vpc) Id() string {
	return v.VpcId
}

func (v Vpc) Type() string {
	return "VPC"
}

func (v Vpc) Delete(region account.Region, account account.Account) error {
	client, err := vpc.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := vpc.CreateDeleteVpcRequest()
	request.VpcId = v.VpcId

	_, err = client.DeleteVpc(request)
	if err != nil {
		return err
	}

	return nil
}
