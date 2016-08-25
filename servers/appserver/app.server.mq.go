package main

import (
	"strings"

	"github.com/arsgo/ars/cluster"
)

//OnMQConsumerCreate  创建JOB服务
func (a *AppServer) OnMQConsumerCreate(task cluster.TaskItem) (path string) {
	path, err := a.clusterClient.CreateMQConsumer(task.Name, a.snap.getMQSnap(task.Name))
	if err != nil {
		a.Log.Error("mq consumer创建失败: ", err)
		return
	}
	a.Log.Infof("::start mq consumer:[%s] %s", task.Name, task.Script)

	return
}

//OnMQConsumerClose 关闭JOB服务
func (a *AppServer) OnMQConsumerClose(task cluster.TaskItem, path string) {
	err := a.clusterClient.CloseNode(path)
	if err != nil {
		return
	}
	return
}

func (a *AppServer) getMQConsumerTask(tasks []cluster.TaskItem) (tks []cluster.TaskItem) {
	if tasks == nil {
		tasks = make([]cluster.TaskItem, 0, 0)
	}
	tks = make([]cluster.TaskItem, 0, len(tasks))
	for _, v := range tasks {
		if strings.EqualFold(v.Type, "mq") &&
			strings.EqualFold(v.Method, "consume") && !v.Disable {
			tks = append(tks, v)
		}
	}
	return
}

func (a *AppServer) bindMQConsumer(tasks []cluster.TaskItem) {
	conusmers := a.getMQConsumerTask(tasks)
	if len(conusmers) == 0 {
		a.Log.Info(" -> 没有可用mq consumer或未配置")
	}
	a.mqService.UpdateTasks(conusmers)
}
