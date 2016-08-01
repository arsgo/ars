package main

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/colinyl/ars/base"
	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/proxy"
	"github.com/colinyl/ars/rpc"
	"github.com/colinyl/ars/server"
	"github.com/colinyl/ars/servers/config"
	"github.com/colinyl/lib4go/concurrent"
	"github.com/colinyl/lib4go/logger"
)

const (
	SERVER_MASTER = "master"
	SERVER_SLAVE  = "slave"
)

//RCServer RC Server
type RCServer struct {
	clusterClient   cluster.IClusterClient
	startSync       base.Sync
	domain          string
	isReconnect     bool
	IsMaster        bool
	currentServices *concurrent.ConcurrentMap
	crossDomain     *concurrent.ConcurrentMap //map[string]cluster.IClusterClient
	crossService    *concurrent.ConcurrentMap //map[string]map[string][]string
	Log             logger.ILogger
	rcRPCServer     *server.RPCServer  //RC Server服务供RPC调用
	spRPCClient     *rpc.RPCClient     //SP Server调用客户端
	rcRPCHandler    server.IRPCHandler //RC Server处理程序
	snap            RCSnap
	loggerName      string
	version         string
}

//NewRCServer 创建RC Server服务器
func NewRCServer() *RCServer {
	rc := &RCServer{loggerName: "rc.server", version: "0.1.10"}
	rc.currentServices = concurrent.NewConcurrentMap()
	rc.crossDomain = concurrent.NewConcurrentMap()
	rc.crossService = concurrent.NewConcurrentMap()
	rc.Log, _ = logger.Get(rc.loggerName)
	rc.startSync = base.NewSync(1)
	return rc
}

//init 初始化服务
func (rc *RCServer) init() (err error) {
	defer rc.recover()
	cfg, err := config.Get()
	if err != nil {
		os.Exit(1)
		return
	}
	rc.Log.Infof(" -> 初始化 %s...", cfg.Domain)
	rc.domain = cfg.Domain
	rc.clusterClient, err = cluster.GetClusterClient(cfg.Domain, cfg.IP, rc.loggerName, cfg.ZKServers...)
	if err != nil {
		return
	}
	rc.spRPCClient = rpc.NewRPCClient(rc.clusterClient, rc.loggerName)
	rc.snap = RCSnap{Domain: cfg.Domain, Server: SERVER_SLAVE, ip: cfg.IP, rcServer: rc, Version: rc.version}
	rc.rcRPCHandler = proxy.NewRPCClientProxy(rc.clusterClient, rc.spRPCClient, rc.snap, rc.loggerName)
	rc.rcRPCServer = server.NewRPCServer(rc.rcRPCHandler, rc.loggerName)
	return nil
}
func (rc *RCServer) recover() {
	if r := recover(); r != nil {
		rc.Log.Fatal(r, string(debug.Stack()))
	}
}

//Start 启动服务
func (rc *RCServer) Start() (err error) {
	defer rc.recover()
	rc.Log.Info(" -> 启动RC Server...")
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
	if rc.BindRCServer() != nil {
		return
	}
	rc.startSync.Wait()
	go rc.startRefreshSnap()
	go rc.startMonitor()
	rc.Log.Info(" -> RC Server 启动完成...")
	return nil
}

//Stop 停止服务
func (rc *RCServer) Stop() error {
	rc.Log.Info(" -> 退出RC Server...")
	defer rc.recover()
	rc.clusterClient.CloseRCServer(rc.snap.Path)
	time.Sleep(time.Millisecond * 10)
	rc.clusterClient.Close()
	rc.spRPCClient.Close()
	rc.rcRPCServer.Stop()
	rc.Log.Info("::rc server closed")
	return nil
}
