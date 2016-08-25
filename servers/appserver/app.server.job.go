package main

import (
	"strings"

	"github.com/arsgo/ars/cluster"
)

func (a *AppServer) bindJobConsumer(tasks []cluster.TaskItem) {
	remoteTasks := a.getJobConsumerTask(tasks)
	a.jobServer.Stop()
	if len(remoteTasks) == 0 {
		a.Log.Info(" -> 没有可用job consumer或未配置")
	}
	a.jobServer.Start()
	a.jobServer.UpdateTasks(remoteTasks)
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

//OnRemoteJobCreate  创建JOB服务
func (a *AppServer) OnRemoteJobCreate(task cluster.TaskItem) (path string) {
	path, err := a.clusterClient.CreateJobConsumer(task.Name, a.snap.getJobConsumerSnap(a.jobServer.Address))
	if err != nil {
		a.Log.Error("job consumer创建失败: ", err)
		a.jobServerCollector.Error(task.Name)
		return
	}
	a.Log.Infof("::start job consumer:[%s] %s", task.Name, task.Script)
	return
}

//OnRemoteJobClose 关闭JOB服务
func (a *AppServer) OnRemoteJobClose(task cluster.TaskItem, path string) {
	a.Log.Error("关闭job:", path)
	err := a.clusterClient.CloseNode(path)
	if err != nil {
		return
	}
	return
}
