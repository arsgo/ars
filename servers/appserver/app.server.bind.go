package main

import (
	"errors"
	"time"

	"github.com/arsgo/ars/cluster"
)

//BindRCServer 绑定RPC调用服务
func (a *AppServer) BindRCServer(configs []*cluster.RCServerItem, err error) (er error) {
	defer a.recover()
	if a.disableRPC {
		return
	}
	er = err
	if err != nil {
		return
	}
	var tasks []string
	for _, v := range configs {
		tasks = append(tasks, v.Address)
	}
	services := make(map[string][]string)
	services["*"] = tasks
	if len(tasks) == 0 {
		a.Log.Error(" -> 没有可用的 rc server")
		er = errors.New(" -> 没有可用的 rc server")
	} else {
		a.Log.Info("::bind rc server ", tasks)
	}
	a.rpcClient.ResetRPCServer(services)
	return
}

//BindTask 绑定本地任务
func (a *AppServer) BindTask(config *cluster.AppServerTask, err error) (er error) {
	defer a.recover()
	if err != nil || config == nil {
		return
	}
	a.snapRefresh = time.Second * time.Duration(config.Config.SnapRefresh)
	a.disableRPC = config.Config.DisableRPC
	a.scriptPool.SetPackages(config.Config.Libs...)
	a.rpcClient.SetPoolSize(config.Config.RPC.MinSize, config.Config.RPC.MaxSize)
	a.scriptPool.SetPoolSize(config.Config.RPC.MinSize, config.Config.RPC.MaxSize)
	a.BindHttpServer(config.Server)
	a.BindLocalJobs(config.LocalJobs)
	a.BindLocalTask(config.Tasks)
	return
}

//OnJobCreate  创建JOB服务
func (a *AppServer) OnJobCreate(task cluster.TaskItem) (path string) {
	path, err := a.clusterClient.CreateJobConsumer(task.Name, a.snap.GetJobSnap(a.jobServer.Address))
	if err != nil {
		a.Log.Error("job consumer创建失败: ", err)
		return
	}
	a.Log.Infof("::start job consumer:[%s] %s", task.Name, task.Script)
	return
}

//OnJobClose 关闭JOB服务
func (a *AppServer) OnJobClose(task cluster.TaskItem, path string) {
	err := a.clusterClient.CloseNode(path)
	if err != nil {
		return
	}
	return
}
