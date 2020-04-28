package aliyun

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type EcsDisks struct{}

type EcsDisk struct {
	ecs.Disk

	DiskType string
}

func init() {
	cloud.RegisterService(EcsDisks{})
}

func (d EcsDisks) IsGlobal() bool {
	return false
}

func (d EcsDisks) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := ecs.CreateDescribeDisksRequest()
	request.PageSize = "99"
	response, err := client.DescribeDisks(request)
	if err != nil {
		return nil, err
	}

	disks := make([]cloud.Resource, 0)
	for _, disk := range response.Disks.Disk {
		disks = append(disks, EcsDisk{Disk: disk, DiskType: disk.Type})
	}

	return disks, nil
}

func (d EcsDisk) Id() string {
	return d.DiskId
}

func (d EcsDisk) Type() string {
	return fmt.Sprintf("Block storage %s disk", d.DiskType)
}

func (d EcsDisk) Delete(region account.Region, account account.Account) error {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := ecs.CreateDeleteDiskRequest()
	request.DiskId = d.DiskId
	_, err = client.DeleteDisk(request)
	if err != nil {
		return err
	}

	return nil
}
