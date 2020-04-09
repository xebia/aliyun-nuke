package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type VpcVpcs struct{}

type VpcVpc struct {
	vpc.Vpc
}

func init() {
	cloud.RegisterService(VpcVpcs{})
}

func (v VpcVpcs) IsGlobal() bool {
	return false
}

func (v VpcVpcs) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
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
		vpcs = append(vpcs, VpcVpc{Vpc: vpcItem})
	}

	return vpcs, nil
}

func (v VpcVpc) Id() string {
	return v.VpcId
}

func (v VpcVpc) Type() string {
	return "VPC"
}

func (v VpcVpc) Delete(region account.Region, account account.Account) error {
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
