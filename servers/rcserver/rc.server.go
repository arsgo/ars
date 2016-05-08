package main

import (
	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/rpcclient"
	"github.com/colinyl/ars/rpcproxy"
	"github.com/colinyl/ars/rpcserver"
	"github.com/colinyl/ars/servers/config"
	"github.com/colinyl/lib4go/logger"
)

const (
	SERVER_MASTER = "master"
	SERVER_SLAVE  = "slave"
)

//RCServer RC Server
type RCServer struct {
	clusterClient     cluster.IClusterClient
	IsMaster          bool
	Log               *logger.Logger
	rcRPCServer       *rpcserver.RPCServer       //RC Server服务供RPC调用
	spRPCClient       *rpcclient.RPCClient       //SP Server调用客户端
	rcRPCProxyHandler rpcserver.RPCScriptHandler //RC Server处理程序
	snap              RCSnap
}

//NewRCServer 创建RC Server服务器
func NewRCServer() *RCServer {
	rc := &RCServer{}
	rc.Log, _ = logger.New("rc server", true)
	return rc
}

//init 初始化服务
func (rc *RCServer) init() (err error) {
	rc.clusterClient, err = cluster.GetClusterClient(config.Get().Domain, config.Get().IP, config.Get().ZKServers...)
	if err != nil {
		return
	}
	rc.snap = RCSnap{Domain: config.Get().Domain, Server: SERVER_SLAVE}
	rc.spRPCClient = rpcclient.NewRPCClient()
	rc.rcRPCProxyHandler = rpcproxy.NewRPCClientProxyHandler(rc.clusterClient, rc.spRPCClient, rc.snap)
	rc.rcRPCServer = rpcserver.NewRPCServer(rc.rcRPCProxyHandler)
	return nil
}

//Start 启动服务
func (rc *RCServer) Start() (err error) {
	if err = rc.init(); err != nil {
		return
	}
	//启动RPC服务,供APP,SP调用
	rc.rcRPCServer.Start()
	//绑定RC服务
	err = rc.BindRCServer()
	if err != nil {
		return
	}
	go rc.StartRefreshSnap()
	return nil
}

//Stop 停止服务
func (rc *RCServer) Stop() error {
	defer func() {
		recover()
	}()
	rc.rcRPCServer.Stop()
	rc.Log.Info("::rc server closed")
	return nil
}
