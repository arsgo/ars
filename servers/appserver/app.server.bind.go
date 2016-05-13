package main

import (
	"fmt"
	"strings"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/httpserver"
	"github.com/colinyl/lib4go/scheduler"
)

//BindRCServer 绑定RPC调用服务
func (a *AppServer) BindRCServer(configs []*cluster.RCServerItem, err error) error {
	var tasks []string
	for _, v := range configs {
		tasks = append(tasks, v.Address)
	}
	services := make(map[string][]string)
	services["*"] = tasks
	a.rpcClient.ResetRPCServer(services)
	return nil
}

//BindTask 绑定本地任务
func (a *AppServer) BindTask(config *cluster.AppServerStartupConfig, err error) (er error) {
	a.ResetAPPSnap()
	scheduler.Stop()
	for _, v := range config.Tasks {
		er = a.scriptPool.Pool.PreLoad(v.Script, 1)
		if er != nil {
			a.Log.Error("load task`s script error in:", v.Script, ",", er)
			continue
		}
		scheduler.AddTask(v.Trigger, scheduler.NewTask(v.Script, func(name interface{}) {
			a.Log.Infof("start:%s", name)
			rtvalues, err := a.scriptPool.Call(name.(string), "{}", "{}")
			if err != nil {
				a.Log.Error(err)
			} else {
				a.Log.Infof("result:%d,%s", len(rtvalues), strings.Join(rtvalues, ","))
			}
		}))
	}
	if len(config.Tasks) > 0 {
		scheduler.Start()
	} else {
		scheduler.Stop()
	}
	if a.httpServer != nil {
		a.httpServer.Stop()
	}
	if config.Server != nil && len(config.Server.Routes) > 0 &&
		strings.EqualFold(strings.ToLower(config.Server.ServerType), "http") {
		a.httpServer, err = httpserver.NewHttpScriptServer(config.Server.Routes, a.scriptPool.Call)
		if err == nil {
			a.httpServer.Start()
			a.snap.Server = fmt.Sprint(a.snap.ip, a.httpServer.Address)
		} else {
			a.Log.Error(err)
		}

	}
	if len(config.Jobs) > 0 {
		a.jobConsumerRPCServer.Stop()
		a.jobConsumerRPCServer.Start()
		a.snap.Address = fmt.Sprint(a.snap.ip, a.jobConsumerRPCServer.Address)
		a.jobConsumerRPCServer.UpdateTasks(config.Jobs)
	} else {
		a.jobConsumerRPCServer.Stop()
	}
	return nil
}

//OnJobCreate  创建JOB服务
func (a *AppServer) OnJobCreate(task cluster.TaskItem) (path string) {
	path, err := a.clusterClient.CreateJobConsumer(task.Name, a.snap.GetSnap())
	if err != nil {
		a.Log.Error(err)
		return
	}
	a.Log.Info("::start job service:", task.Name)
	return
}

//OnJobClose 关闭JOB服务
func (app *AppServer) OnJobClose(task cluster.TaskItem, path string) {
	err := app.clusterClient.CloseJobConsumer(path)
	if err != nil {
		return
	}
	return
}
