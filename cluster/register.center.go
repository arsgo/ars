package cluster

import (
	"encoding/json"
	"time"

	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

const (
	SERVER_MASTER    = "master"
	SERVER_SLAVE     = "slave"
	rcServerRoot     = "@domain/rc/servers"
	rcServerPath     = "@domain/rc/servers/rc_"
	rcServerNodePath = "@domain/rc/servers/@name"
	//rcServerValue    = `{"domain":"@domain","path":"@path","ip":"@ip","port":"@port","server":"@type","online":"@online","lastPublish":"@pst","last":"@last"}`
	rcServerConfig = "@domain/rc/config"

	//jobRoot             = "@domain/job"
	jobConfigPath       = "@domain/job/config"
	jobConsumerRoot     = "@domain/job/servers/@jobName"
	jobConsumerRealPath = "@domain/job/servers/@jobName/@path"
)

type rcSnap struct {
	Domain  string          `json:"domain"`
	Path    string          `json:"path"`
	Address string          `json:"address"`
	Server  string          `json:"server"`
	Last    string           `json:"last"`
	Sys     *sysMonitorInfo `json:"sys"`
}

func (a rcSnap) GetSnap() string {
	snap := a
	snap.Last =  time.Now().Format("20060102150405")
	snap.Sys, _ = GetSysMonitorInfo()
	buffer, _ := json.Marshal(&snap)
	return string(buffer)
}

//-------------------------register center----------------------------
type rcServer struct {
	Path           string
	//IP             string
	//Port           string
	Server         string
	dataMap        utility.DataMap
	IsMasterServer bool	
	jobCallback        func(config *JobConfigs, err error)
	Log                *logger.Logger
	rpcServer          *rpcservice.RPCServer
	rcServerRoot       string
	rcServerPath       string
	servicePublishPath string
	serviceRoot        string
	jobConfigPath      string
	spServerPool       *rpcservice.RPCServerPool
	spServicesMap      *servicesMap
	zkClient           *clusterClient
	snap               rcSnap
}

//JobConfigItem job config item
type JobConfigItem struct {
	Name        string
	Script      string
	Trigger     string
	Concurrency int
}
type JobConsumerValue struct {
	Address string
}

//JobConfigs job configs
type JobConfigs struct {
	Jobs map[string]JobConfigItem
}

func NewRCServer() *rcServer {
	rc := &rcServer{}
	rc.Log, _ = logger.New("rc server", true)
	return rc
}
func (rc *rcServer) init() error {
	rc.zkClient = NewClusterClient()
	rc.dataMap = rc.zkClient.dataMap.Copy()
	rc.rcServerRoot = rc.dataMap.Translate(rcServerRoot)
	rc.servicePublishPath = rc.dataMap.Translate(servicePublishPath)
	rc.serviceRoot = rc.dataMap.Translate(serviceRoot)
	rc.jobConfigPath = rc.dataMap.Translate(jobConfigPath)
	rc.rcServerPath = rc.dataMap.Translate(rcServerPath)
	rc.spServerPool = rpcservice.NewRPCServerPool()
	rc.spServicesMap = NewServiceMap()
	rc.snap = rcSnap{Domain: rc.zkClient.Domain, Server: SERVER_SLAVE}
	return nil
}

func (r *rcServer) Start() (err error) {
	if err = r.init(); err != nil {
		return
	}
	r.StartRPCServer()
	err = r.Bind()
	if err != nil {
		return
	}
	r.WatchJobChange(func(config *JobConfigs, err error) {
		r.BindScheduler(config, err)
	})
	r.WatchServiceChange(func(services map[string][]string, err error) {
		r.BindSPServer(services)
	})
	go r.StartRefreshSnap()
	return nil
}

func (r *rcServer) Stop() error {
	defer func() {
		recover()
	}()
	r.zkClient.ZkCli.Close()
	if r.rpcServer != nil {
		r.rpcServer.Stop()
	}
	r.Log.Info("::rc server closed")
	return nil
}
