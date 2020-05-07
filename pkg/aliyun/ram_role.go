package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
	"strings"
)

type RamRoles struct{}

type RamRole struct {
	ram.RoleInListRoles

	Policies []ram.PolicyInListPoliciesForRole
}

func init() {
	cloud.RegisterService(RamRoles{})
}

func (r RamRoles) IsGlobal() bool {
	return true
}

func (r RamRoles) List(region account.Region, account account.Account, force bool) ([]cloud.Resource, error) {
	client, err := ram.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := ram.CreateListRolesRequest()
	request.Scheme = "https"
	response, err := client.ListRoles(request)
	if err != nil {
		return nil, err
	}

	roles := make([]cloud.Resource, 0)
	for _, role := range response.Roles.Role {
		if force || !strings.HasPrefix(role.RoleName, "Aliyun") {
			policies, err := fetchPoliciesForRole(client, role.RoleName)
			if err != nil {
				return nil, err
			}

			roles = append(roles, RamRole{
				RoleInListRoles: role,
				Policies:        policies,
			})
		}
	}

	return roles, nil
}

func (r RamRole) Id() string {
	return r.RoleName
}

func (r RamRole) Type() string {
	return "RAM role"
}

func (r RamRole) Delete(region account.Region, account account.Account) error {
	client, err := ram.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	// Detach policies from user
	for _, policy := range r.Policies {
		request := ram.CreateDetachPolicyFromRoleRequest()
		request.Scheme = "https"
		request.PolicyName = policy.PolicyName
		request.PolicyType = policy.PolicyType
		request.RoleName = r.RoleName
		_, err := client.DetachPolicyFromRole(request)
		if err != nil {
			return err
		}
	}

	// Delete user
	request := ram.CreateDeleteRoleRequest()
	request.Scheme = "https"
	request.RoleName = r.RoleName

	_, err = client.DeleteRole(request)
	if err != nil {
		return err
	}

	return nil
}

func fetchPoliciesForRole(client *ram.Client, roleName string) ([]ram.PolicyInListPoliciesForRole, error) {
	request := ram.CreateListPoliciesForRoleRequest()
	request.Scheme = "https"
	request.RoleName = roleName
	response, err := client.ListPoliciesForRole(request)
	if err != nil {
		return nil, err
	}

	return response.Policies.Policy, nil
}
