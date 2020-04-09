package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"time"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type VpcNatGateways struct{}

type VpcNatGateway struct {
	vpc.NatGateway

	SnatTables map[string][]vpc.SnatTableEntry
}

func init() {
	cloud.RegisterService(VpcNatGateways{})
}

func (n VpcNatGateways) IsGlobal() bool {
	return false
}

// List returns a list of all NAT gateways
func (n VpcNatGateways) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := vpc.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := vpc.CreateDescribeNatGatewaysRequest()
	request.PageSize = "50"
	response, err := client.DescribeNatGateways(request)
	if err != nil {
		return nil, err
	}

	natGateways := make([]cloud.Resource, 0)
	for _, natGatewayItem := range response.NatGateways.NatGateway {
		snatTables, err := fetchSnatTables(client, natGatewayItem.SnatTableIds)
		if err != nil {
			return nil, err
		}
		natGateways = append(natGateways, VpcNatGateway{NatGateway: natGatewayItem, SnatTables: snatTables})
	}

	return natGateways, nil
}

func (n VpcNatGateway) Id() string {
	return n.NatGatewayId
}

func (n VpcNatGateway) Type() string {
	return "NAT gateway"
}

func (n VpcNatGateway) Delete(region account.Region, account account.Account) error {
	client, err := vpc.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	err = deleteNatGatewayAndAwait(client, n)
	if err != nil {
		return err
	}

	return nil
}

func fetchSnatTables(client *vpc.Client, ids vpc.SnatTableIdsInDescribeNatGateways) (map[string][]vpc.SnatTableEntry, error) {
	snatTables := make(map[string][]vpc.SnatTableEntry)
	for _, id := range ids.SnatTableId {
		request := vpc.CreateDescribeSnatTableEntriesRequest()
		request.SnatTableId = id
		response, err := client.DescribeSnatTableEntries(request)
		if err != nil {
			return nil, err
		}
		snatTables[id] = response.SnatTableEntries.SnatTableEntry
	}
	return snatTables, nil
}

func deleteNatGatewayAndAwait(client *vpc.Client, n VpcNatGateway) error {
	deleteRequest := vpc.CreateDeleteNatGatewayRequest()
	deleteRequest.NatGatewayId = n.NatGatewayId
	deleteRequest.Force = "true"
	_, err := client.DeleteNatGateway(deleteRequest)
	if err != nil {
		return err
	}

	isDeleted := false
	for !isDeleted {
		getRequest := vpc.CreateDescribeNatGatewaysRequest()
		getRequest.PageSize = "50"
		response, err := client.DescribeNatGateways(getRequest)
		if err != nil {
			return err
		}

		found := false
		for _, natGateway := range response.NatGateways.NatGateway {
			if natGateway.NatGatewayId == n.NatGatewayId {
				found = true
			}
		}

		if !found {
			isDeleted = true
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
