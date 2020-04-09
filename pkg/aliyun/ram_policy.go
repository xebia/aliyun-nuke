package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type RamPolicies struct{}

type RamPolicy struct {
	ram.PolicyInListPolicies
}

func init() {
	cloud.RegisterService(RamPolicies{})
}

func (p RamPolicies) IsGlobal() bool {
	return true
}

func (p RamPolicies) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := ram.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := ram.CreateListPoliciesRequest()
	request.Scheme = "https"
	request.PolicyType = "Custom"
	response, err := client.ListPolicies(request)
	if err != nil {
		return nil, err
	}

	policies := make([]cloud.Resource, 0)
	for _, policy := range response.Policies.Policy {
		policies = append(policies, RamPolicy{
			PolicyInListPolicies: policy,
		})
	}

	return policies, nil
}

func (p RamPolicy) Id() string {
	return p.PolicyName
}

func (p RamPolicy) Type() string {
	return "RAM policy"
}

func (p RamPolicy) Delete(region account.Region, account account.Account) error {
	client, err := ram.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := ram.CreateDeletePolicyRequest()
	request.Scheme = "https"
	request.PolicyName = p.PolicyName

	_, err = client.DeletePolicy(request)
	if err != nil {
		return err
	}

	return nil
}
