package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type RamUsers struct{}

type RamUser struct {
	ram.UserInListUsers

	AccessKeys []account.Credentials
	Policies   []ram.PolicyInListPoliciesForUser
	Groups     []ram.GroupInListGroupsForUser
}

func init() {
	cloud.RegisterService(RamUsers{})
}

func (u RamUsers) IsGlobal() bool {
	return true
}

func (u RamUsers) List(region account.Region, account account.Account, force bool) ([]cloud.Resource, error) {
	client, err := ram.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := ram.CreateListUsersRequest()
	request.Scheme = "https"
	response, err := client.ListUsers(request)
	if err != nil {
		return nil, err
	}

	users := make([]cloud.Resource, 0)
	for _, user := range response.Users.User {
		accessKeys, err := fetchAccessKeysForUser(client, user.UserName)
		if err != nil {
			return nil, err
		}

		policies, err := fetchPoliciesForUser(client, user.UserName)
		if err != nil {
			return nil, err
		}

		groups, err := fetchGroupsForUser(client, user.UserName)
		if err != nil {
			return nil, err
		}

		users = append(users, RamUser{
			UserInListUsers: user,
			AccessKeys:      accessKeys,
			Policies:        policies,
			Groups:          groups,
		})
	}

	return users, nil
}

func (u RamUser) Id() string {
	return u.UserName
}

func (u RamUser) Type() string {
	return "RAM user"
}

func (u RamUser) Delete(region account.Region, account account.Account) error {
	client, err := ram.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	// Delete access keys
	for _, accessKey := range u.AccessKeys {
		request := ram.CreateDeleteAccessKeyRequest()
		request.Scheme = "https"
		request.UserAccessKeyId = accessKey.AccessKeyID
		request.UserName = u.UserName
		_, err := client.DeleteAccessKey(request)
		if err != nil {
			return err
		}
	}

	// Detach policies from user
	for _, policy := range u.Policies {
		request := ram.CreateDetachPolicyFromUserRequest()
		request.Scheme = "https"
		request.PolicyName = policy.PolicyName
		request.PolicyType = policy.PolicyType
		request.UserName = u.UserName
		_, err := client.DetachPolicyFromUser(request)
		if err != nil {
			return err
		}
	}

	// Remove user from groups
	for _, group := range u.Groups {
		request := ram.CreateRemoveUserFromGroupRequest()
		request.Scheme = "https"
		request.GroupName = group.GroupName
		request.UserName = u.UserName
		_, err := client.RemoveUserFromGroup(request)
		if err != nil {
			return err
		}
	}

	// Delete user
	request := ram.CreateDeleteUserRequest()
	request.Scheme = "https"
	request.UserName = u.UserName

	_, err = client.DeleteUser(request)
	if err != nil {
		return err
	}

	return nil
}

func fetchAccessKeysForUser(client *ram.Client, username string) ([]account.Credentials, error) {
	request := ram.CreateListAccessKeysRequest()
	request.Scheme = "https"
	request.UserName = username
	response, err := client.ListAccessKeys(request)
	if err != nil {
		return nil, err
	}

	accessKeys := make([]account.Credentials, 0)
	for _, accessKey := range response.AccessKeys.AccessKey {
		accessKeys = append(accessKeys, account.Credentials{
			AccessKeyID: accessKey.AccessKeyId,
		})
	}
	return accessKeys, nil
}

func fetchPoliciesForUser(client *ram.Client, username string) ([]ram.PolicyInListPoliciesForUser, error) {
	request := ram.CreateListPoliciesForUserRequest()
	request.Scheme = "https"
	request.UserName = username
	response, err := client.ListPoliciesForUser(request)
	if err != nil {
		return nil, err
	}

	return response.Policies.Policy, nil
}

func fetchGroupsForUser(client *ram.Client, username string) ([]ram.GroupInListGroupsForUser, error) {
	request := ram.CreateListGroupsForUserRequest()
	request.Scheme = "https"
	request.UserName = username
	response, err := client.ListGroupsForUser(request)
	if err != nil {
		return nil, err
	}

	return response.Groups.Group, nil
}
