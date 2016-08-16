package main

import (
	"fmt"
	"time"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/ars/mq"
	"github.com/arsgo/ars/proxy"
	"github.com/arsgo/ars/rpc"
	"github.com/arsgo/ars/script"
	"github.com/arsgo/ars/server"
	"github.com/arsgo/ars/servers/config"
	"github.com/arsgo/lib4go/concurrent"
	"github.com/arsgo/lib4go/logger"
)

var (
	eModeShared = "share"
	eModeAlone  = "alone"
)

//SPServer SPServer
type SPServer struct {
	Log                 logger.ILogger
	snapLogger          logger.ILogger
	startSync           base.Sync
	domain              string
	mode                string
	serviceConfig       string
	timerReloadRCServer *base.TimerCall
	mqService           *mq.MQConsumerService
	rpcClient           *rpc.RPCClient
	rpcServer           *server.RPCServer  //RPC 服务器
	rpcScriptProxy      *proxy.ScriptProxy //RPC Server 脚本处理程序
	clusterClient       cluster.IClusterClient
	scriptPool          *script.ScriptPool //脚本引擎池
	dbPool              *concurrent.ConcurrentMap
	snap                SPSnap
	loggerName          string
	version             string
	ip                  string
}

//NewSPServer 创建SP server服务器
func NewSPServer() (sp *SPServer, err error) {
	sp = &SPServer{loggerName: "sp.server", version: "0.1.10"}
	sp.startSync = base.NewSync(2)
	sp.timerReloadRCServer = base.NewTimerCall(time.Second*5, time.Microsecond, sp.reloadRCServer)
	sp.Log, err = logger.Get(sp.loggerName)
	if err != nil {
		return
	}
	sp.snapLogger, err = logger.Get("sp.snap")
	if err != nil {
		return
	}
	sp.snapLogger.Show(false)
	sp.dbPool = concurrent.NewConcurrentMap()
	return
}

//init 初始化服务器
func (sp *SPServer) init() (err error) {
	defer sp.recover()
	cfg, err := config.Get()
	if err != nil {
		return
	}
	sp.Log.Infof(" -> 初始化 %s...", cfg.Domain)
	sp.ip = cfg.IP
	sp.clusterClient, err = cluster.NewDomainClusterClient(cfg.Domain, cfg.IP, sp.loggerName, cfg.ZKServers...)
	if err != nil {
		return
	}
	sp.domain = cfg.Domain
	sp.snap = SPSnap{ip: cfg.IP, Version: sp.version}
	sp.rpcClient = rpc.NewRPCClient(sp.clusterClient, sp.loggerName)
	sp.scriptPool, err = script.NewScriptPool(sp.clusterClient, sp.rpcClient, sp.GetScriptBinder(), sp.loggerName)
	if err != nil {
		return
	}
	sp.rpcScriptProxy = proxy.NewScriptProxy(sp.clusterClient, sp.scriptPool, sp.loggerName)
	sp.rpcScriptProxy.OnOpenTask = sp.OnSPServiceCreate
	sp.rpcScriptProxy.OnCloseTask = sp.OnSPServiceClose
	sp.rpcServer = server.NewRPCServer(sp.rpcScriptProxy, sp.loggerName, sp.collectReporter)
	sp.mqService, err = mq.NewMQConsumerService(sp.clusterClient, mq.NewMQScriptHandler(sp.scriptPool, sp.loggerName), sp.loggerName)
	return
}

//Start 启动SP Server服务器
func (sp *SPServer) Start() (err error) {
	defer sp.recover()
	sp.Log.Info(" -> 启动 sp server...")
	if err = sp.init(); err != nil {
		sp.Log.Error(err)
		return
	}
	if !sp.clusterClient.WaitForConnected() {
		return
	}
	sp.rpcServer.Start()
	sp.snap.Address = fmt.Sprint(sp.snap.ip, sp.rpcServer.Address)
	sp.clusterClient.WatchSPTaskChange(func() {
		defer sp.startSync.Done("INIT.TASK.BIND")
		sp.rebindService()
	})
	sp.clusterClient.WatchRCServerChange(func(config []*cluster.RCServerItem, err error) {
		defer sp.startSync.Done("INIT.RCSRV.BIND")
		sp.BindRCServer(config, err)
	})
	sp.startSync.Wait()
	go sp.startRefreshSnap()
	go sp.startMonitor()
	sp.Log.Info(" -> sp server 启动完成...")
	return nil
}

//Stop 停止SP Server服务器
func (sp *SPServer) Stop() error {
	defer sp.recover()
	sp.Log.Info(" -> 退出 sp server...")
	sp.clusterClient.Close()
	sp.rpcClient.Close()
	sp.rpcServer.Stop()
	return nil
}
