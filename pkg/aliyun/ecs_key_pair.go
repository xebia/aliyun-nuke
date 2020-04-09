package aliyun

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type EcsKeyPairs struct{}

type EcsKeyPair struct {
	ecs.KeyPair
}

func init() {
	cloud.RegisterService(EcsKeyPairs{})
}

func (k EcsKeyPairs) IsGlobal() bool {
	return false
}

func (k EcsKeyPairs) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := ecs.CreateDescribeKeyPairsRequest()
	request.PageSize = "50"
	response, err := client.DescribeKeyPairs(request)
	if err != nil {
		return nil, err
	}

	keyPairs := make([]cloud.Resource, 0)
	for _, keyPair := range response.KeyPairs.KeyPair {
		keyPairs = append(keyPairs, EcsKeyPair{KeyPair: keyPair})
	}

	return keyPairs, nil
}

func (k EcsKeyPair) Id() string {
	return k.KeyPairName
}

func (k EcsKeyPair) Type() string {
	return "SSH key pair"
}

func (k EcsKeyPair) Delete(region account.Region, account account.Account) error {
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
