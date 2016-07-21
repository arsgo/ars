package main

import "github.com/colinyl/ars/cluster"

//BindRCServer 绑定RPC调用服务
func (sp *SPServer) BindRCServer(configs []*cluster.RCServerItem, err error) error {
	var tasks []string
	for _, v := range configs {
		tasks = append(tasks, v.Address)
	}
	services := make(map[string][]string)
	services["*"] = tasks
	sp.rpcClient.ResetRPCServer(services)
	return nil
}

//rebindService 重新绑定SP所有服务列表
func (sp *SPServer) rebindService() {
	sp.Log.Info(" -> start bind services")
	task, err := sp.clusterClient.GetServiceTasks()
	if err != nil {
		sp.Log.Error(err)
		return
	}
	sp.Log.Info("rpc pool size min:", task.Config.RPC.MinSize, ",max:", task.Config.RPC.MaxSize)
	sp.scriptPool.SetPackages(task.Config.Libs...)
	sp.rpcClient.SetPoolSize(task.Config.RPC.MinSize, task.Config.RPC.MaxSize)
	sp.rpcServer.UpdateTasks(task.Tasks)
	err = sp.mqService.UpdateTasks(task.Tasks)
	if err != nil {
		sp.Log.Error(err)
		return
	}
}

//OnSPServiceCreate 服务创建时同时创建集群节点
func (sp *SPServer) OnSPServiceCreate(task cluster.TaskItem) (path string) {
	sp.Log.Info("::load script:", task.Script, ",minSize:", task.MinSize, ",maxSize:", task.MaxSize)
	sp.scriptPool.PreLoad(task.Script, task.MinSize, task.MaxSize)
	path, err := sp.clusterClient.CreateServiceProvider(task.Name, sp.rpcServer.Address,
		sp.snap.GetSnap(task.Name))
	if err != nil {
		return
	}
	sp.Log.Info("::start service provider:", task.Name)
	return
}

//OnSPServiceClose 服务停止时同时删除集群节点
func (sp *SPServer) OnSPServiceClose(task cluster.TaskItem, path string) {
	err := sp.clusterClient.CloseServiceProvider(path)
	if err != nil {
		return
	}
	return
}
