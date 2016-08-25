package main

import (
	"strings"

	"github.com/arsgo/ars/cluster"
)

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

//BindLocalTask 绑定本地任务，包括MQ Consumer,Job Consumer
func (a *AppServer) BindLocalTask(tasks []cluster.TaskItem) {
	a.bindMQConsumer(tasks)
	a.bindJobConsumer(tasks)
	a.bindLocalJobs(tasks)
}
