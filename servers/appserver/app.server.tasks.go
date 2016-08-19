package main

import (
	"strings"

	"github.com/arsgo/ars/cluster"
)

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
func (a *AppServer) getJobLocalsTask(tasks []cluster.TaskItem) (tks []cluster.TaskItem) {
	if tasks == nil {
		tasks = make([]cluster.TaskItem, 0, 0)
	}
	tks = make([]cluster.TaskItem, 0, len(tasks))
	for _, v := range tasks {
		if strings.EqualFold(v.Type, "job") &&
			strings.EqualFold(v.Method, "local") && !v.Disable {
			tks = append(tks, v)
		}
	}
	return
}
func (a *AppServer) getJobConsumerTask(tasks []cluster.TaskItem) (tks []cluster.TaskItem) {
	if tasks == nil {
		tasks = make([]cluster.TaskItem, 0, 0)
	}
	tks = make([]cluster.TaskItem, 0, len(tasks))
	for _, v := range tasks {
		if strings.EqualFold(v.Type, "job") &&
			strings.EqualFold(v.Method, "consume") && !v.Disable {
			tks = append(tks, v)
		}
	}
	return
}
func (a *AppServer) bindMQConsumer(tasks []cluster.TaskItem) {
	conusmers := a.getMQConsumerTask(tasks)
	if len(conusmers) == 0 {
		a.Log.Info("没有可用mq consumer或未配置")
	}
	a.mqService.UpdateTasks(conusmers)
}

func (a *AppServer) bindJobConsumer(tasks []cluster.TaskItem) {
	remoteTasks := a.getJobConsumerTask(tasks)
	a.jobServer.Stop()
	if len(remoteTasks) == 0 {
		a.Log.Info("没有可用job consumer或未配置")
		return
	}
	a.jobServer.Start()
	a.jobServer.UpdateTasks(remoteTasks)
}

//BindLocalTask 绑定本地任务，包括MQ Consumer,Job Consumer
func (a *AppServer) BindLocalTask(tasks []cluster.TaskItem) {
	a.bindMQConsumer(tasks)
	a.bindJobConsumer(tasks)
	a.bindLocalJobs(tasks)
}
