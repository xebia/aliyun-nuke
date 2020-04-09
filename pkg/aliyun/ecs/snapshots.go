package ecs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type Snapshots struct{}

type Snapshot struct {
	ecs.Snapshot
}

func init() {
	cloud.RegisterService(Snapshots{})
}

func (s Snapshots) IsGlobal() bool {
	return false
}

func (s Snapshots) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := ecs.CreateDescribeSnapshotsRequest()
	request.PageSize = "99"
	response, err := client.DescribeSnapshots(request)
	if err != nil {
		return nil, err
	}

	snapshots := make([]cloud.Resource, 0)
	for _, snapshot := range response.Snapshots.Snapshot {
		snapshots = append(snapshots, Snapshot{Snapshot: snapshot})
	}

	return snapshots, nil
}

func (s Snapshot) Id() string {
	return s.SnapshotId
}

func (s Snapshot) Type() string {
	return "Snapshot"
}

func (s Snapshot) Delete(region account.Region, account account.Account) error {
	client, err := ecs.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := ecs.CreateDeleteSnapshotRequest()
	request.SnapshotId = s.SnapshotId
	_, err = client.DeleteSnapshot(request)
	if err != nil {
		return err
	}

	return nil
}
