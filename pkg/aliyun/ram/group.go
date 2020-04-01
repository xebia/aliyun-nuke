package ram

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type Groups struct{}

type Group struct {
	ram.Group

	Policies []ram.Policy
}

func (g Groups) IsGlobal() bool {
	return true
}

func (g Groups) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := ram.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := ram.CreateListGroupsRequest()
	request.Scheme = "https"
	response, err := client.ListGroups(request)
	if err != nil {
		return nil, err
	}

	groups := make([]cloud.Resource, 0)
	for _, group := range response.Groups.Group {
		policies, err := fetchPoliciesForGroup(client, group.GroupName)
		if err != nil {
			return nil, err
		}

		groups = append(groups, Group{
			Group:    group,
			Policies: policies,
		})
	}

	return groups, nil
}

func (g Group) Id() string {
	return g.GroupName
}

func (g Group) Type() string {
	return "RAM group"
}

func (g Group) Delete(region account.Region, account account.Account) error {
	client, err := ram.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	// Detach policies from user
	for _, policy := range g.Policies {
		request := ram.CreateDetachPolicyFromGroupRequest()
		request.Scheme = "https"
		request.PolicyName = policy.PolicyName
		request.PolicyType = policy.PolicyType
		request.GroupName = g.GroupName
		_, err := client.DetachPolicyFromGroup(request)
		if err != nil {
			return err
		}
	}

	// Delete user
	request := ram.CreateDeleteGroupRequest()
	request.Scheme = "https"
	request.GroupName = g.GroupName

	_, err = client.DeleteGroup(request)
	if err != nil {
		return err
	}

	return nil
}

func fetchPoliciesForGroup(client *ram.Client, groupName string) ([]ram.Policy, error) {
	request := ram.CreateListPoliciesForGroupRequest()
	request.Scheme = "https"
	request.GroupName = groupName
	response, err := client.ListPoliciesForGroup(request)
	if err != nil {
		return nil, err
	}

	return response.Policies.Policy, nil
}
