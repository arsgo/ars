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
			strings.EqualFold(v.Method, "consume") {
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
			strings.EqualFold(v.Method, "consume") {
			tks = append(tks, v)
		}
	}
	return
}
func (a *AppServer) bindMQConsumer(tasks []cluster.TaskItem) {
	conusmers := a.getMQConsumerTask(tasks)
	a.mqService.UpdateTasks(conusmers)
}
func (a *AppServer) bindJobConsumer(tasks []cluster.TaskItem) {
	conusmers := a.getJobConsumerTask(tasks)
	a.jobServer.Stop()
	if len(conusmers) == 0 {
		return
	}
	a.jobServer.Start()
	a.jobServer.UpdateTasks(conusmers)

}

//BindLocalTask 绑定本地任务，包括MQ Consumer,Job Consumer
func (a *AppServer) BindLocalTask(tasks []cluster.TaskItem) {
	a.bindMQConsumer(tasks)
	a.bindJobConsumer(tasks)
}
