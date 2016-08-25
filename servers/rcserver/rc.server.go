package main

import (
	"runtime/debug"
	"time"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/ars/proxy"
	"github.com/arsgo/ars/rpc"
	"github.com/arsgo/ars/server"
	"github.com/arsgo/ars/servers/config"
	"github.com/arsgo/lib4go/concurrent"
	"github.com/arsgo/lib4go/logger"
)

//RCServer RC Server
type RCServer struct {
	clusterClient        cluster.IClusterClient
	startSync            base.Sync
	clusterServers       []string
	isReconnect          bool
	IsMaster             bool
	rpcServerCollector   *base.Collector
	schedulerCollector   *base.Collector
	currentServices      *concurrent.ConcurrentMap
	crossDomain          *concurrent.ConcurrentMap //map[string]cluster.IClusterClient
	crossServices        *concurrent.ConcurrentMap //map[string]map[string][]string
	timerRebindServices  *base.TimerCall
	timerPublishServices *base.TimerCall
	Log                  logger.ILogger
	snapLogger           logger.ILogger
	rcRPCServer          *server.RPCServer  //RC Server服务供RPC调用
	spRPCClient          *rpc.RPCClient     //SP Server调用客户端
	rcRPCHandler         server.IRPCHandler //RC Server处理程序
	snap                 RCSnap
	conf                 *config.SysConfig
	loggerName           string
	version              string
}

//NewRCServer 创建RC Server服务器
func NewRCServer(conf *config.SysConfig) (rc *RCServer, err error) {
	rc = &RCServer{loggerName: "rc.server", version: "0.1.15", conf: conf}
	rc.currentServices = concurrent.NewConcurrentMap()
	rc.crossDomain = concurrent.NewConcurrentMap()
	rc.crossServices = concurrent.NewConcurrentMap()
	rc.timerRebindServices = base.NewTimerCall(time.Second*5, time.Millisecond, rc.rebindLocalServices)
	rc.timerPublishServices = base.NewTimerCall(time.Second*3, time.Second, rc.PublishNow)
	rc.rpcServerCollector = base.NewCollector()
	rc.schedulerCollector = base.NewCollector()
	rc.startSync = base.NewSync(1)
	rc.Log, err = logger.Get(rc.loggerName)
	if err != nil {
		return
	}
	rc.snapLogger, err = logger.Get("rc.snap")
	if err != nil {
		return
	}
	//rc.snapLogger.Show(false)
	return
}

//init 初始化服务
func (rc *RCServer) init() (err error) {
	defer rc.recover()
	rc.Log.Infof(" -> 初始化 %s...", rc.conf.Domain)
	rc.clusterServers = rc.conf.ZKServers
	rc.clusterClient, err = cluster.NewDomainClusterClient(rc.conf.Domain, rc.conf.IP, rc.loggerName, rc.conf.ZKServers...)
	if err != nil {
		return
	}
	rc.spRPCClient = rpc.NewRPCClient(rc.clusterClient, rc.loggerName)
	rc.snap = RCSnap{Domain: rc.conf.Domain, Server: cluster.SERVER_UNKNOWN, rcServer: rc, Version: rc.version, Refresh: 60}
	rc.rcRPCHandler = proxy.NewRPCClientProxy(rc.clusterClient, rc.spRPCClient, rc.loggerName)
	rc.rcRPCServer = server.NewRPCServer(rc.rcRPCHandler, rc.loggerName, rc.rpcServerCollector)
	return
}

//Start 启动服务
func (rc *RCServer) Start() (err error) {
	defer rc.recover()
	rc.Log.Info(" -> 启动 rc server...")
	if err = rc.init(); err != nil {
		rc.Log.Error(err)
		return
	}

	if !rc.clusterClient.WaitForConnected() {
		return
	}
	//启动RPC服务,供APP,SP调用
	rc.rcRPCServer.Start()

	//绑定RC服务
	if err = rc.BindRCServer(); err != nil {
		return
	}
	rc.startSync.Wait()
	go rc.startRefreshSnap()
	go rc.startMonitor()
	go rc.clearMem()
	rc.Log.Info(" -> rc server 启动完成...")
	return
}

//Stop 停止服务
func (rc *RCServer) Stop() error {
	rc.Log.Info(" -> 退出 rc server...")
	defer rc.recover()
	rc.spRPCClient.Close()
	rc.rcRPCServer.Stop()
	rc.clusterClient.CloseRCServer(rc.snap.path)
	rc.clusterClient.Close()
	cross := rc.crossDomain.GetAll()
	for _, v := range cross {
		cls := v.(cluster.IClusterClient)
		cls.Close()
	}
	return nil
}
func (rc *RCServer) recover() (err error) {
	if r := recover(); r != nil {
		rc.Log.Fatal(r, string(debug.Stack()))
	}
	return
}
