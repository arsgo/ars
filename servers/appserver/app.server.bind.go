package main

import "github.com/arsgo/ars/cluster"

//BindRCServer 绑定RPC调用服务
func (a *AppServer) BindRCServer(configs []*cluster.RCServerItem, err error) error {
	defer a.recover()
	var tasks []string
	for _, v := range configs {
		tasks = append(tasks, v.Address)
	}
	services := make(map[string][]string)
	services["*"] = tasks
	a.Log.Info(" -> bind rc server (", len(services), ")")
	a.rpcClient.ResetRPCServer(services)
	return nil
}

//BindTask 绑定本地任务
func (a *AppServer) BindTask(config *cluster.AppServerStartupConfig, err error) (er error) {
	defer a.recover()
	if config == nil {
		return
	}
	a.Log.Info("rpc pool size min:", config.Config.RPC.MinSize, ",max:", config.Config.RPC.MaxSize)
	a.scriptPool.SetPackages(config.Config.Libs...)
	a.rpcClient.SetPoolSize(config.Config.RPC.MinSize, config.Config.RPC.MaxSize)
	a.scriptPool.SetPoolSize(config.Config.RPC.MinSize, config.Config.RPC.MaxSize)
	a.ResetAPPSnap()
	a.BindHttpServer(config.Server)
	a.BindLocalJobs(config.LocalJobs)
	a.BindLocalTask(config.Tasks)
	return
}

//OnJobCreate  创建JOB服务
func (a *AppServer) OnJobCreate(task cluster.TaskItem) (path string) {
	path, err := a.clusterClient.CreateJobConsumer(task.Name, a.snap.GetJobSnap(a.jobServer.Address))
	if err != nil {
		a.Log.Error(err)
		return
	}
	a.Log.Infof("::start job consumer:[%s] %s", task.Name, task.Script)
	return
}

//OnJobClose 关闭JOB服务
func (a *AppServer) OnJobClose(task cluster.TaskItem, path string) {
	err := a.clusterClient.CloseJobConsumer(path)
	if err != nil {
		return
	}
	return
}
