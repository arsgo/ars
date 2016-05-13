package main

import (
	"fmt"
	"sync"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/rpcproxy"
	"github.com/colinyl/ars/servers/config"
	"github.com/colinyl/lib4go/logger"
)

var (
	eModeShared = "share"
	eModeAlone  = "alone"
)

//SPServer SPServer
type SPServer struct {
	Log            *logger.Logger
	lk             sync.Mutex
	mode           string
	serviceConfig  string
	rpcClient      *rpcproxy.RPCClient
	rpcServer      *rpcproxy.RPCServer        //RPC 服务器
	rpcScriptProxy *rpcproxy.RPCScriptHandler //RPC Server 脚本处理程序
	clusterClient  cluster.IClusterClient
	scriptPool     *rpcproxy.ScriptPool //脚本引擎池
	snap           SPSnap
}

//NewSPServer 创建SP server服务器
func NewSPServer() *SPServer {
	sp := &SPServer{}
	sp.Log, _ = logger.New("sp server", true)
	return sp
}

//init 初始化服务器
func (sp *SPServer) init() (err error) {
	cfg := config.Get()
	sp.clusterClient, err = cluster.GetClusterClient(cfg.Domain, cfg.IP, cfg.ZKServers...)
	if err != nil {
		return
	}
	sp.snap = SPSnap{ip: cfg.IP}
	sp.rpcClient = rpcproxy.NewRPCClient()
	sp.scriptPool, err = rpcproxy.NewScriptPool(sp.clusterClient, sp.rpcClient)
	if err != nil {
		return
	}
	sp.rpcScriptProxy = rpcproxy.NewRPCScriptHandler(sp.clusterClient, sp.scriptPool)
	sp.rpcScriptProxy.OnOpenTask = sp.OnSPServiceCreate
	sp.rpcScriptProxy.OnCloseTask = sp.OnSPServiceClose
	sp.rpcServer = rpcproxy.NewRPCServer(sp.rpcScriptProxy)
	return
}

//Start 启动SP Server服务器
func (sp *SPServer) Start() (err error) {
	if err = sp.init(); err != nil {
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
	defer recover()
	sp.clusterClient.Close()
	sp.rpcServer.Stop()
	sp.Log.Info("::sp server closed")
	return nil
}
