package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type EssScheduledTasks struct{}

type EssScheduledTask struct {
	ess.ScheduledTask
}

func init() {
	cloud.RegisterService(EssScheduledTasks{})
}

func (s EssScheduledTasks) IsGlobal() bool {
	return false
}

func (s EssScheduledTasks) List(region account.Region, account account.Account, force bool) ([]cloud.Resource, error) {
	client, err := ess.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	request := ess.CreateDescribeScheduledTasksRequest()
	request.PageSize = "50"
	response, err := client.DescribeScheduledTasks(request)
	if err != nil {
		return nil, err
	}

	scheduledTasks := make([]cloud.Resource, 0)
	for _, scheduledTask := range response.ScheduledTasks.ScheduledTask {
		scheduledTasks = append(scheduledTasks, EssScheduledTask{ScheduledTask: scheduledTask})
	}

	return scheduledTasks, nil
}

func (s EssScheduledTask) Id() string {
	return s.ScheduledTaskId
}

func (s EssScheduledTask) Type() string {
	return "ESS scheduled task"
}

func (s EssScheduledTask) Delete(region account.Region, account account.Account) error {
	client, err := ess.NewClientWithAccessKey(string(region), account.AccessKeyID, account.AccessKeySecret)
	if err != nil {
		return err
	}

	request := ess.CreateDeleteScheduledTaskRequest()
	request.ScheduledTaskId = s.ScheduledTaskId

	_, err = client.DeleteScheduledTask(request)
	if err != nil {
		return err
	}

	return nil
}
