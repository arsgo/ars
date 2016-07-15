package main

import (
	"sync"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/httpserver"
	"github.com/colinyl/ars/mqservice"
	"github.com/colinyl/ars/rpcproxy"
	"github.com/colinyl/ars/servers/config"
	"github.com/colinyl/lib4go/logger"
)

//AppServer app server服务器
type AppServer struct {
	JobAddress               map[string]string
	domain	string
	Log                      logger.ILogger
	clusterClient            cluster.IClusterClient
	jobConsumerScriptHandler *rpcproxy.RPCScriptHandler //本地JOB Consumer提供的RPC接口,使用的代理处理程序为脚本处理
	jobConsumerRPCServer     *rpcproxy.RPCServer        //接收JOB事件调用,改事件将触发脚本执行
	rpcClient                *rpcproxy.RPCClient        //RPC远程调用客户端,调用RC Server提供的RPC服务
	scriptPool               *rpcproxy.ScriptPool       //脚本池,用于缓存JOB Consumer脚本和本地task任务执行脚本
	lk                       sync.Mutex
	httpServer               *httpserver.HTTPScriptServer
	mqService                *mqservice.MQConsumerService
	snap                     AppSnap
	loggerName               string
}

//NewAPPServer 创建APP Server服务器
func NewAPPServer() *AppServer {
	app := &AppServer{loggerName: "app.server"}
	app.JobAddress = make(map[string]string)
	app.Log, _ = logger.Get(app.loggerName, true)
	return app
}

//init 初始化服务器
func (app *AppServer) init() (err error) {
	defer app.recover()
	cfg, err := config.Get()
	if err != nil {
		return
	}
	app.Log.Infof(" -> 初始化 %s...",cfg.Domain)	
	app.clusterClient, err = cluster.GetClusterClient(cfg.Domain, cfg.IP, app.loggerName, cfg.ZKServers...)
	if err != nil {
		return
	}
	app.domain=cfg.Domain
	app.rpcClient = rpcproxy.NewRPCClient(app.clusterClient, app.loggerName)
	app.scriptPool, err = rpcproxy.NewScriptPool(app.clusterClient, app.rpcClient, make(map[string]interface{}), app.loggerName)
	app.jobConsumerScriptHandler = rpcproxy.NewRPCScriptHandler(app.clusterClient, app.scriptPool, app.loggerName)
	app.jobConsumerScriptHandler.OnOpenTask = app.OnJobCreate
	app.jobConsumerScriptHandler.OnCloseTask = app.OnJobClose
	app.jobConsumerRPCServer = rpcproxy.NewRPCServer(app.jobConsumerScriptHandler, app.loggerName)
	app.mqService, err = mqservice.NewMQConsumerService(app.clusterClient, mqservice.NewMQScriptHandler(app.scriptPool, app.loggerName), app.loggerName)
	app.snap = AppSnap{ip: cfg.IP, appserver: app}
	app.snap.Address = cfg.IP
	return
}

//Start 启动服务器
func (app *AppServer) Start() (err error) {
	defer app.recover()
	app.Log.Info(" -> 启动APP Server...")
	if err = app.init(); err != nil {
		app.Log.Error(err)
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
	defer app.recover()
	app.Log.Info(" -> 退出AppServer...")
	app.clusterClient.Close()
	app.rpcClient.Close()
	app.scriptPool.Close()
	app.jobConsumerRPCServer.Stop()
	if app.httpServer != nil {
		app.httpServer.Stop()
	}

	app.Log.Info("::app server closed")
	return nil
}
