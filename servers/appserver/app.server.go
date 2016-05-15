package main

import (
	"sync"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/httpserver"
	"github.com/colinyl/ars/rpcproxy"
	"github.com/colinyl/ars/servers/config"
	"github.com/colinyl/lib4go/logger"
)

//AppServer app server服务器
type AppServer struct {
	JobAddress               map[string]string
	Log                      *logger.Logger
	clusterClient            cluster.IClusterClient
	jobConsumerScriptHandler *rpcproxy.RPCScriptHandler //本地JOB Consumer提供的RPC接口,使用的代理处理程序为脚本处理
	jobConsumerRPCServer     *rpcproxy.RPCServer        //接收JOB事件调用,改事件将触发脚本执行
	rpcClient                *rpcproxy.RPCClient        //RPC远程调用客户端,调用RC Server提供的RPC服务
	scriptPool               *rpcproxy.ScriptPool       //脚本池,用于缓存JOB Consumer脚本和本地task任务执行脚本
	lk                       sync.Mutex
	httpServer               *httpserver.HttpScriptServer
	snap                     AppSnap
}

//NewAPPServer 创建APP Server服务器
func NewAPPServer() *AppServer {
	app := &AppServer{}
	app.JobAddress = make(map[string]string)
	app.Log, _ = logger.New("app server", true)
	return app
}

//init 初始化服务器
func (app *AppServer) init() (err error) {
	app.clusterClient, err = cluster.GetClusterClient(config.Get().Domain, config.Get().IP, config.Get().ZKServers...)
	if err != nil {
		return
	}
	app.snap = AppSnap{ip: config.Get().IP}
	app.rpcClient = rpcproxy.NewRPCClient(app.clusterClient)
	app.snap.Address = config.Get().IP
	app.scriptPool, err = rpcproxy.NewScriptPool(app.clusterClient, app.rpcClient)
	app.jobConsumerScriptHandler = rpcproxy.NewRPCScriptHandler(app.clusterClient, app.scriptPool)
	app.jobConsumerScriptHandler.OnOpenTask = app.OnJobCreate
	app.jobConsumerScriptHandler.OnCloseTask = app.OnJobClose
	app.jobConsumerRPCServer = rpcproxy.NewRPCServer(app.jobConsumerScriptHandler)

	return
}

//Start 启动服务器
func (app *AppServer) Start() (err error) {
	if err = app.init(); err != nil {
		return
	}
	app.clusterClient.WatchRCServerChange(func(config []*cluster.RCServerItem, err error) {
		app.BindRCServer(config, err)
	})

	app.clusterClient.WatchAppTaskChange(func(config *cluster.AppServerStartupConfig, err error) error {
		app.BindTask(config, err)
		return nil
	})
	go app.StartRefreshSnap()
	return nil
}

//Stop 停止服务器
func (app *AppServer) Stop() error {
	defer func() {
		recover()
	}()
	app.clusterClient.Close()
	app.jobConsumerRPCServer.Stop()
	if app.httpServer != nil {
		app.httpServer.Stop()
	}

	app.Log.Info("::app server closed")
	return nil
}
