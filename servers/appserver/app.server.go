package main

import (
	"time"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/ars/mq"
	"github.com/arsgo/ars/proxy"
	"github.com/arsgo/ars/rpc"
	"github.com/arsgo/ars/script"
	"github.com/arsgo/ars/server"
	"github.com/arsgo/ars/servers/config"
	"github.com/arsgo/lib4go/logger"
)

//AppServer app server服务器
type AppServer struct {
	JobAddress          map[string]string
	domain              string
	startSync           base.Sync
	Log                 logger.ILogger
	clusterClient       cluster.IClusterClient
	timerReloadRCServer *base.TimerCall
	disableRPC          bool
	scriptPorxy         *proxy.ScriptProxy //本地脚本处理
	jobServer           *server.RPCServer  //接收JOB事件调用,改事件将触发脚本执行
	rpcClient           *rpc.RPCClient     //RPC远程调用客户端,调用RC Server提供的RPC服务
	scriptPool          *script.ScriptPool //脚本池,用于缓存JOB Consumer脚本和本地task任务执行脚本
	httpServer          *server.HTTPScriptServer
	mqService           *mq.MQConsumerService
	snapRefresh         time.Duration
	snap                AppSnap
	loggerName          string
	conf                *config.SysConfig
	version             string
}

//NewAPPServer 创建APP Server服务器
func NewAPPServer() (app *AppServer, err error) {
	app = &AppServer{loggerName: "app.server", version: "0.1.10"}
	app.timerReloadRCServer = base.NewTimerCall(time.Second*5, time.Microsecond, app.reloadRCServer)
	app.startSync = base.NewSync(2)
	app.JobAddress = make(map[string]string)
	app.Log, err = logger.Get(app.loggerName)
	if err != nil {
		return
	}
	app.conf, err = config.Get()
	if err != nil {
		app.Log.Error(err)
		return
	}
	return
}

//init 初始化服务器
func (app *AppServer) init() (err error) {
	defer app.recover()
	app.Log.Infof(" -> 初始化 %s...", app.conf.Domain)
	app.clusterClient, err = cluster.NewDomainClusterClient(app.conf.Domain, app.conf.IP, app.loggerName, app.conf.ZKServers...)
	if err != nil {
		return
	}
	app.domain = app.conf.Domain
	app.rpcClient = rpc.NewRPCClient(app.clusterClient, app.loggerName)
	app.scriptPool, err = script.NewScriptPool(app.clusterClient, app.rpcClient, make(map[string]interface{}), app.loggerName)
	if err != nil {
		return
	}
	app.scriptPorxy = proxy.NewScriptProxy(app.clusterClient, app.scriptPool, app.loggerName)
	app.scriptPorxy.OnOpenTask = app.OnJobCreate
	app.scriptPorxy.OnCloseTask = app.OnJobClose
	app.jobServer = server.NewRPCServer(app.scriptPorxy, app.loggerName, app.collectReporter)
	app.mqService, err = mq.NewMQConsumerService(app.clusterClient, mq.NewMQScriptHandler(app.scriptPool, app.loggerName), app.loggerName)
	if err != nil {
		return
	}
	app.snap = AppSnap{ip: app.conf.IP, appserver: app, Version: app.version}
	app.snap.Address = app.conf.IP
	return
}

//Start 启动服务器
func (app *AppServer) Start() (err error) {
	defer app.recover()

	app.Log.Info(" -> 启动 app server...")
	if err = app.init(); err != nil {
		app.Log.Error(err)
		return
	}
	if !app.clusterClient.WaitForConnected() {
		return
	}

	app.clusterClient.WatchAppTaskChange(func(config *cluster.AppServerTask, err error) error {
		defer app.startSync.Done("INIT.BIND.TASK")
		app.BindTask(config, err)
		return nil
	})
	app.clusterClient.WatchRCServerChange(func(config []*cluster.RCServerItem, err error) {
		defer app.startSync.Done("INIT.BIND.RCSRV")
		app.BindRCServer(config, err)
	})

	app.startSync.Wait()
	go app.startMonitor()
	go app.StartRefreshSnap()

	app.Log.Info(" -> app server 启动完成...")
	return nil
}

//Stop 停止服务器
func (app *AppServer) Stop() error {
	defer app.recover()
	app.Log.Info(" -> 退出 app server...")
	app.clusterClient.Close()
	app.rpcClient.Close()
	app.scriptPool.Close()
	app.jobServer.Stop()
	if app.httpServer != nil {
		app.httpServer.Stop()
	}
	return nil
}
