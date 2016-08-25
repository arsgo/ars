package main

import (
	"errors"

	"github.com/arsgo/ars/cluster"
)

//BindRCServer 绑定RPC调用服务
func (sp *SPServer) BindRCServer(configs []*cluster.RCServerItem, err error) (er error) {
	defer sp.recover()
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
		sp.Log.Error(" -> 没有可用的 rc server")
		er = errors.New(" -> 没有可用的 rc server")
	} else {
		sp.Log.Info("::bind rc server ", tasks)
	}
	sp.rpcClient.ResetRPCServer(services)
	return nil
}

//rebindService 重新绑定SP所有服务列表
func (sp *SPServer) rebindService() {
	task, err := sp.clusterClient.GetSPServerTask(sp.conf.IP)
	if err != nil {
		sp.Log.Error(err)
		return
	}
	sp.scriptPool.SetPackages(task.Config.Libs...)
	sp.rpcClient.SetPoolSize(task.Config.RPC.MinSize, task.Config.RPC.MaxSize)

	if sp.rpcServer.UpdateTasks(task.Tasks) > 0 {
		sp.Log.Info(" -> 本地服务已更新:", len(task.Tasks))
		sp.snapLogger.Infof("--------------------services-----------------\n\t\t\t\t\t  %+v\n\t\t\t\t\t  ----------------------------------------------",
			sp.rpcServer.GetServices())
	}
	err = sp.mqService.UpdateTasks(task.Tasks)
	if err != nil {
		sp.Log.Error(err)
		return
	}
}

//OnSPServiceCreate 服务创建时同时创建集群节点
func (sp *SPServer) OnSPServiceCreate(task cluster.TaskItem) (path string) {
	sp.scriptPool.PreLoad(task.Script, task.MinSize, task.MaxSize)
	path, err := sp.clusterClient.CreateSPServer(task.Name, sp.rpcServer.Address,
		sp.snap.getDefSnap(task.Name))
	if err != nil {
		sp.Log.Errorf("创建sp server node error:%v", err)
		return
	}
	sp.Log.Info("::start service:", task.Name)
	return
}

//OnSPServiceClose 服务停止时同时删除集群节点
func (sp *SPServer) OnSPServiceClose(task cluster.TaskItem, path string) {
	err := sp.clusterClient.CloseSPServer(path)
	if err != nil {
		sp.Log.Errorf("关闭服务失败:%s,%v", path, err)
		return
	}
	return
}
