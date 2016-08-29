package main

import (
	"errors"
	"time"

	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/ars/snap"
	"github.com/arsgo/lib4go/utility"
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
	a.snap.Refresh = utility.GetMax2(config.Config.SnapRefresh, 120, 60)
	if config.Config.SnapRefresh > 0 && config.Config.SnapRefresh < 60 {
		a.Log.Error(" -> 快照刷新时间不能低于60秒")
	}
	snap.ResetTicker(time.Second * time.Duration(a.snap.Refresh))
	a.disableRPC = config.Config.DisableRPC
	a.scriptPool.SetPackages(config.Config.Libs...)
	a.rpcClient.SetPoolSize(config.Config.RPC.MinSize, config.Config.RPC.MaxSize)
	a.scriptPool.SetPoolSize(config.Config.RPC.MinSize, config.Config.RPC.MaxSize)
	a.BindAPIServer(config.Server)
	a.BindLocalTask(config.Tasks)

	return
}
