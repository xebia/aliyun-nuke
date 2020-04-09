package ram

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type Policies struct{}

type Policy struct {
	ram.PolicyInListPolicies
}

func init() {
	cloud.RegisterService(Policies{})
}

func (p Policies) IsGlobal() bool {
	return true
}

func (p Policies) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
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
		policies = append(policies, Policy{
			PolicyInListPolicies: policy,
		})
	}

	return policies, nil
}

func (p Policy) Id() string {
	return p.PolicyName
}

func (p Policy) Type() string {
	return "RAM policy"
}

func (p Policy) Delete(region account.Region, account account.Account) error {
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
