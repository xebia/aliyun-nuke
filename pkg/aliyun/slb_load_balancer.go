package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type SlbLoadBalancers struct{}

type SlbLoadBalancer struct {
	slb.LoadBalancer
}

func init() {
	cloud.RegisterService(SlbLoadBalancers{})
}

func (l SlbLoadBalancers) IsGlobal() bool {
	return false
}

func (l SlbLoadBalancers) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := slb.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := slb.CreateDescribeLoadBalancersRequest()
	request.PageSize = "50"
	response, err := client.DescribeLoadBalancers(request)
	if err != nil {
		return nil, err
	}

	loadBalancers := make([]cloud.Resource, 0)
	for _, loadBalancer := range response.LoadBalancers.LoadBalancer {
		loadBalancers = append(loadBalancers, SlbLoadBalancer{LoadBalancer: loadBalancer})
	}

	return loadBalancers, nil
}

func (l SlbLoadBalancer) Id() string {
	return l.LoadBalancerId
}

func (l SlbLoadBalancer) Type() string {
	return "SLB Load balancer"
}

func (l SlbLoadBalancer) Delete(region account.Region, account account.Account) error {
	client, err := slb.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := slb.CreateDeleteLoadBalancerRequest()
	request.LoadBalancerId = l.LoadBalancerId

	_, err = client.DeleteLoadBalancer(request)
	if err != nil {
		return err
	}

	return nil
}
