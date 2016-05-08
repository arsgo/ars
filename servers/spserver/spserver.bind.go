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
	tasks, err := sp.clusterClient.GeServiceTasks()
	if err != nil {
		return
	}
	sp.rpcServer.UpdateTasks(tasks)
}

//OnSPServiceCreate 服务创建时同时创建集群节点
func (sp *SPServer) OnSPServiceCreate(task cluster.TaskItem) (path string) {
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
