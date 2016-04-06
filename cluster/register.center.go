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
	rcServerConfig   = "@domain/configs/rc/config"

	jobRoot          = "@domain/job"
	jobConfigPath    = "@domain/configs/job/config"
	jobConsumerRoot  = "@domain/job/@jobName"
	jobConsumerPath  = "@domain/job/@jobName/job"
	jobConsumerValue = `{"ip":"@ip",last":@now}`
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
}

//JobConfigItem job config item
type JobConfigItem struct {
	Name        string
	Script      string
	Trigger     string
	Concurrency int
}

//JobConfigs job configs
type JobConfigs struct {
	Jobs map[string]JobConfigItem
}

func NewRCServer() *rcServer {
	var err error
	rc := &rcServer{}
	rc.Log, err = logger.New("register center", true)
	rc.dataMap = utility.NewDataMap()
	rc.dataMap.Set("domain", zkClient.Domain)
	rc.dataMap.Set("ip", zkClient.LocalIP)
	rc.dataMap.Set("type", "salve")
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
