package main

import (
	"fmt"
	"runtime/debug"
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
	mode                string
	serviceConfig       string
	rpcServerCollector  *base.Collector
	mqConsumerCollector *base.Collector
	collectorMap        map[string]base.ICollector
	timerReloadRCServer *base.TimerCall
	mqService           *mq.MQConsumerService
	rpcClient           *rpc.RPCClient
	rpcServer           *server.RPCServer  //RPC 服务器
	rpcScriptProxy      *proxy.ScriptProxy //RPC Server 脚本处理程序
	clusterClient       cluster.IClusterClient
	scriptPool          *script.ScriptPool //脚本引擎池
	dbPool              *concurrent.ConcurrentMap
	snap                SPSnap
	conf                *config.SysConfig
	loggerName          string
	version             string
}

//NewSPServer 创建SP server服务器
func NewSPServer(conf *config.SysConfig) (sp *SPServer, err error) {
	sp = &SPServer{loggerName: "sp.server", version: "0.1.10", conf: conf}
	sp.rpcServerCollector = base.NewCollector()
	sp.mqConsumerCollector = base.NewCollector()
	sp.startSync = base.NewSync(2)

	sp.collectorMap = make(map[string]base.ICollector)
	sp.collectorMap[base.TN_SERVICE_PROVIDER] = sp.rpcServerCollector
	sp.collectorMap[base.TN_MQ_CONSUMER] = sp.mqConsumerCollector

	sp.timerReloadRCServer = base.NewTimerCall(time.Second*5, time.Microsecond, sp.reloadRCServer)
	sp.Log, err = logger.Get(sp.loggerName)
	if err != nil {
		return
	}
	sp.snapLogger, err = logger.Get("sp.snap")
	if err != nil {
		return
	}
	//sp.snapLogger.Show(false)
	sp.dbPool = concurrent.NewConcurrentMap()
	logger.MainLoggerName = sp.loggerName
	return
}

//init 初始化服务器
func (sp *SPServer) init() (err error) {
	defer sp.recover()
	sp.Log.Infof(" -> 初始化 %s...", sp.conf.Domain)
	sp.clusterClient, err = cluster.NewDomainClusterClient(sp.conf.Domain, sp.conf.IP, sp.loggerName, sp.conf.ZKServers...)
	if err != nil {
		return
	}
	sp.snap = SPSnap{Version: sp.version, spserver: sp, Refresh: 60}
	sp.rpcClient = rpc.NewRPCClient(sp.clusterClient, sp.loggerName)
	sp.scriptPool, err = script.NewScriptPool(sp.clusterClient, sp.rpcClient, sp.getdbTypeBinder(), sp.loggerName, sp.collectorMap)
	if err != nil {
		return
	}
	sp.rpcScriptProxy = proxy.NewScriptProxy(sp.clusterClient, sp.scriptPool, base.TN_SERVICE_PROVIDER, sp.loggerName)
	sp.rpcScriptProxy.OnOpenTask = sp.OnSPServiceCreate
	sp.rpcScriptProxy.OnCloseTask = sp.OnSPServiceClose
	sp.rpcServer = server.NewRPCServer(sp.rpcScriptProxy, sp.loggerName, sp.rpcServerCollector)
	handler := mq.NewMQScriptHandler(sp.scriptPool, sp.loggerName, sp.OnMQConsumerCreate, sp.OnMQConsumerClose, sp.mqConsumerCollector)
	sp.mqService, err = mq.NewMQConsumerService(sp.clusterClient, handler, sp.loggerName, sp.mqConsumerCollector)
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
	sp.snap.Address = fmt.Sprint(sp.conf.IP, sp.rpcServer.Address)
	sp.clusterClient.WatchSPTaskChange(func(task cluster.SPServerTask, err error) {
		defer sp.startSync.Done("INIT.TASK.BIND")
		sp.bindServiceTask(task, err)
	})
	sp.clusterClient.WatchRCServerChange(func(config []*cluster.RCServerItem, err error) {
		defer sp.startSync.Done("INIT.RCSRV.BIND")
		sp.BindRCServer(config, err)
	})
	sp.startSync.Wait()
	go sp.startRefreshSnap()
	go sp.startMonitor()
	go sp.clearMem()
	sp.Log.Info(" -> sp server 启动完成...")
	return nil
}

//Stop 停止SP Server服务器
func (sp *SPServer) Stop() error {
	defer sp.recover()
	sp.Log.Info(" -> 退出 sp server...")
	sp.clusterClient.Close()
	sp.CloseDB()
	sp.rpcClient.Close()
	sp.rpcServer.Stop()
	return nil
}

func (sp *SPServer) recover() {
	if r := recover(); r != nil {
		sp.Log.Fatal(r, string(debug.Stack()))
	}
}
