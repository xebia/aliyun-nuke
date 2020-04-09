package ecs

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/xebia/aliyun-nuke/pkg/nuker"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type KeyPairs struct{}

type KeyPair struct {
	ecs.KeyPair
}

func init() {
	cloud.RegisterService(KeyPairs{})
}

func (k KeyPairs) IsGlobal() bool {
	return false
}

func (k KeyPairs) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := ecs.CreateDescribeKeyPairsRequest()
	request.PageSize = "99"
	response, err := client.DescribeKeyPairs(request)
	if err != nil {
		return nil, err
	}

	keyPairs := make([]cloud.Resource, 0)
	for _, keyPair := range response.KeyPairs.KeyPair {
		keyPairs = append(keyPairs, KeyPair{KeyPair: keyPair})
	}

	return keyPairs, nil
}

func (k KeyPair) Id() string {
	return k.KeyPairName
}

func (k KeyPair) Type() string {
	return "SSH key pair"
}

func (k KeyPair) Delete(region account.Region, account account.Account) error {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := ecs.CreateDeleteKeyPairsRequest()
	request.KeyPairNames = fmt.Sprintf("[\"%s\"]", k.KeyPairName)
	_, err = client.DeleteKeyPairs(request)
	if err != nil {
		return err
	}

	return nil
}
