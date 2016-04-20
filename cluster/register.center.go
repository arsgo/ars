package cluster

import (
	"log"

	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

const (
	rcServerRoot     = "@domain/rc/servers"
	rcServerPath     = "@domain/rc/servers/rc_"
	rcServerNodePath = "@domain/rc/servers/@name"
	rcServerValue    = `{"domain":"@domain","path":"@path","ip":"@ip","port":"@port","server":"@type","online":"@online","lastPublish":"@pst","last":"@last"}`
	rcServerConfig   = "@domain/rc/config"

	//jobRoot             = "@domain/job"
	jobConfigPath       = "@domain/job/config"
	jobConsumerRoot     = "@domain/job/servers/@jobName"
	jobConsumerRealPath = "@domain/job/servers/@jobName/@path"
)

//-------------------------register center----------------------------
type rcServer struct {
	Path               string
	IP                 string
	Port               string
	Server             string
	dataMap            *utility.DataMap
	IsMasterServer     bool
	Last               int64
	OnlineTime         int64
	LastPublish        int64
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
	zkClient           *zkClientObj
}

//JobConfigItem job config item
type JobConfigItem struct {
	Name        string
	Script      string
	Trigger     string
	Concurrency int
}
type JobConsumerValue struct {
	IP string
}

//JobConfigs job configs
type JobConfigs struct {
	Jobs map[string]JobConfigItem
}

func NewRCServer() *rcServer {
	var err error
	rc := &rcServer{}
	rc.Log, err = logger.New("rc server", true)
	rc.zkClient = NewZKClient()
	rc.dataMap = rc.zkClient.dataMap.Copy()
	rc.dataMap.Set("type", "slave")
	rc.rcServerRoot = rc.dataMap.Translate(rcServerRoot)
	rc.servicePublishPath = rc.dataMap.Translate(servicePublishPath)
	rc.serviceRoot = rc.dataMap.Translate(serviceRoot)
	rc.jobConfigPath = rc.dataMap.Translate(jobConfigPath)
	rc.rcServerPath = rc.dataMap.Translate(rcServerPath)
	rc.spServerPool = rpcservice.NewRPCServerPool()
	rc.spServicesMap = NewServiceMap()
	if err != nil {
		log.Print(err)
	}
	return rc
}
func (r *rcServer) Close() {
	defer func() {
		recover()
	}()
	r.zkClient.ZkCli.Close()
	if r.rpcServer != nil {
		r.rpcServer.Stop()
	}
	r.Log.Info("::rc server closed")
}
