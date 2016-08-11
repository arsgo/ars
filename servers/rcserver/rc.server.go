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

const (
	SERVER_MASTER = "master"
	SERVER_SLAVE  = "slave"
)

//RCServer RC Server
type RCServer struct {
	clusterClient        cluster.IClusterClient
	startSync            base.Sync
	clusterServers       []string
	domain               string
	isReconnect          bool
	IsMaster             bool
	currentServices      *concurrent.ConcurrentMap
	crossDomain          *concurrent.ConcurrentMap //map[string]cluster.IClusterClient
	crossService         *concurrent.ConcurrentMap //map[string]map[string][]string
	timerRebindServices  *base.TimerCall
	timerPublishServices *base.TimerCall
	Log                  logger.ILogger
	snapLogger           logger.ILogger
	rcRPCServer          *server.RPCServer  //RC Server服务供RPC调用
	spRPCClient          *rpc.RPCClient     //SP Server调用客户端
	rcRPCHandler         server.IRPCHandler //RC Server处理程序
	snap                 RCSnap
	loggerName           string
	version              string
}

//NewRCServer 创建RC Server服务器
func NewRCServer() (rc *RCServer, err error) {
	rc = &RCServer{loggerName: "rc.server", version: "0.1.10"}
	rc.currentServices = concurrent.NewConcurrentMap()
	rc.crossDomain = concurrent.NewConcurrentMap()
	rc.crossService = concurrent.NewConcurrentMap()
	rc.timerRebindServices = base.NewTimerCall(time.Second*5, time.Millisecond, rc.rebindLocalServices)
	rc.timerPublishServices = base.NewTimerCall(time.Second*3, time.Second, rc.PublishNow)
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
	cfg, err := config.Get()
	if err != nil {
		return
	}
	rc.Log.Infof(" -> 初始化 %s...", cfg.Domain)
	rc.domain = cfg.Domain
	rc.clusterServers = cfg.ZKServers
	rc.clusterClient, err = cluster.NewDomainClusterClient(cfg.Domain, cfg.IP, rc.loggerName, cfg.ZKServers...)
	if err != nil {
		return
	}
	rc.spRPCClient = rpc.NewRPCClient(rc.clusterClient, rc.loggerName)
	rc.snap = RCSnap{Domain: cfg.Domain, Server: SERVER_SLAVE, ip: cfg.IP, rcServer: rc, Version: rc.version}
	rc.rcRPCHandler = proxy.NewRPCClientProxy(rc.clusterClient, rc.spRPCClient, rc.snap, rc.loggerName)
	rc.rcRPCServer = server.NewRPCServer(rc.rcRPCHandler, rc.loggerName, rc.collectReporter)
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
	rc.Log.Info(" -> rc server 启动完成...")
	return
}

//Stop 停止服务
func (rc *RCServer) Stop() error {
	rc.Log.Info(" -> 退出 rc server...")
	defer rc.recover()
	rc.spRPCClient.Close()
	rc.rcRPCServer.Stop()
	rc.clusterClient.CloseRCServer(rc.snap.Path)
	rc.clusterClient.Close()
	return nil
}
func (rc *RCServer) recover() {
	if r := recover(); r != nil {
		rc.Log.Fatal(r, string(debug.Stack()))
	}
}
