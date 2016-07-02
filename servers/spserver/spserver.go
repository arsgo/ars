package main

import (
	"fmt"
	"sync"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/mqservice"
	"github.com/colinyl/ars/rpcproxy"
	"github.com/colinyl/ars/servers/config"
	"github.com/colinyl/lib4go/concurrent"
	"github.com/colinyl/lib4go/logger"
)

var (
	eModeShared = "share"
	eModeAlone  = "alone"
)

//SPServer SPServer
type SPServer struct {
	Log            logger.ILogger
	lk             sync.Mutex
	mode           string
	serviceConfig  string
	mqService      *mqservice.MQConsumerService
	rpcClient      *rpcproxy.RPCClient
	rpcServer      *rpcproxy.RPCServer        //RPC 服务器
	rpcScriptProxy *rpcproxy.RPCScriptHandler //RPC Server 脚本处理程序
	clusterClient  cluster.IClusterClient
	scriptPool     *rpcproxy.ScriptPool //脚本引擎池
	dbPool         concurrent.ConcurrentMap
	snap           SPSnap
	loggerName     string
}

//NewSPServer 创建SP server服务器
func NewSPServer() *SPServer {
	sp := &SPServer{loggerName: "sp.server"}
	sp.Log, _ = logger.Get(sp.loggerName, true)
	return sp
}

//init 初始化服务器
func (sp *SPServer) init() (err error) {
	defer sp.recover()
	sp.Log.Info(" -> 初始化SP Server...")
	cfg, err := config.Get()
	if err != nil {
		return
	}
	sp.clusterClient, err = cluster.GetClusterClient(cfg.Domain, cfg.IP,sp.loggerName, cfg.ZKServers...)
	if err != nil {
		return
	}
	sp.snap = SPSnap{ip: cfg.IP}
	sp.rpcClient = rpcproxy.NewRPCClient(sp.clusterClient,sp.loggerName)
	sp.scriptPool, err = rpcproxy.NewScriptPool(sp.clusterClient, sp.rpcClient, sp.GetScriptBinder(),sp.loggerName)
	if err != nil {
		return
	}
	sp.rpcScriptProxy = rpcproxy.NewRPCScriptHandler(sp.clusterClient, sp.scriptPool,sp.loggerName)
	sp.rpcScriptProxy.OnOpenTask = sp.OnSPServiceCreate
	sp.rpcScriptProxy.OnCloseTask = sp.OnSPServiceClose
	sp.rpcServer = rpcproxy.NewRPCServer(sp.rpcScriptProxy,sp.loggerName)
	sp.mqService, err = mqservice.NewMQConsumerService(sp.clusterClient, mqservice.NewMQScriptHandler(sp.scriptPool,sp.loggerName),sp.loggerName)
	sp.dbPool = concurrent.NewConcurrentMap()
	return
}

//Start 启动SP Server服务器
func (sp *SPServer) Start() (err error) {
	defer sp.recover()
	sp.Log.Info(" -> 启动SP Server...")
	if err = sp.init(); err != nil {
		sp.Log.Error(err)
		return
	}

	sp.rpcServer.Start()
	sp.snap.Address = fmt.Sprint(sp.snap.ip, sp.rpcServer.Address)
	sp.clusterClient.WatchSPTaskChange(func() {
		sp.rebindService()
	})
	sp.clusterClient.WatchRCServerChange(func(config []*cluster.RCServerItem, err error) {
		sp.BindRCServer(config, err)
	})

	go sp.StartRefreshSnap()
	return nil
}

//Stop 停止SP Server服务器
func (sp *SPServer) Stop() error {
	defer sp.recover()
	sp.Log.Info(" -> 退出SP Server...")
	sp.clusterClient.Close()
	sp.rpcClient.Close()
	sp.rpcServer.Stop()
	sp.Log.Info("::sp server closed")
	return nil
}
