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
	domain                   string
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
	conf                     *config.SysConfig
}

//NewAPPServer 创建APP Server服务器
func NewAPPServer() (app *AppServer, err error) {
	app = &AppServer{loggerName: "app.server"}
	app.JobAddress = make(map[string]string)
	app.Log, err = logger.Get(app.loggerName, true)
	if err != nil {
		return
	}
	app.conf, err = config.Get()
	if err != nil {
		return
	}
	return
}

//init 初始化服务器
func (app *AppServer) init() (err error) {
	defer app.recover()

	app.Log.Infof(" -> 初始化 %s...", app.conf.Domain)
	app.clusterClient, err = cluster.GetClusterClient(app.conf.Domain, app.conf.IP, app.loggerName, app.conf.ZKServers...)
	if err != nil {
		return
	}
	app.domain = app.conf.Domain
	app.rpcClient = rpcproxy.NewRPCClient(app.clusterClient, app.loggerName)
	app.scriptPool, err = rpcproxy.NewScriptPool(app.clusterClient, app.rpcClient, make(map[string]interface{}), app.loggerName)
	app.jobConsumerScriptHandler = rpcproxy.NewRPCScriptHandler(app.clusterClient, app.scriptPool, app.loggerName)
	app.jobConsumerScriptHandler.OnOpenTask = app.OnJobCreate
	app.jobConsumerScriptHandler.OnCloseTask = app.OnJobClose
	app.jobConsumerRPCServer = rpcproxy.NewRPCServer(app.jobConsumerScriptHandler, app.loggerName)
	app.mqService, err = mqservice.NewMQConsumerService(app.clusterClient, mqservice.NewMQScriptHandler(app.scriptPool, app.loggerName), app.loggerName)
	app.snap = AppSnap{ip: app.conf.IP, appserver: app}
	app.snap.Address = app.conf.IP
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
	if !app.clusterClient.WatchConnected() {
		return
	}
	app.clusterClient.WatchAppTaskChange(func(config *cluster.AppServerStartupConfig, err error) error {
		app.BindTask(config, err)
		return nil
	})
	app.clusterClient.WatchRCServerChange(func(config []*cluster.RCServerItem, err error) {
		app.BindRCServer(config, err)
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
