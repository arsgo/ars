package cluster

import (
	"log"

	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

const (
	appServerConfig   = "@domain/configs/app/@ip"
	appServerRoot     = "@domain/app/servers"
	appServerRootPath = "@domain/app/servers/@ip"
)

type AutoConfig struct {
	Trigger string
	Script  string
	Input   string
}

type AppConfig struct {
	Status string
	Auto   []*AutoConfig
}
type RCServerConfig struct {
	Domain string
	IP     string
	Port   string
	Server string
	Online string
}

type appServer struct {
	dataMap         *utility.DataMap
	Last            int64
	Log             *logger.Logger
	appServerConfig string
	rcServerRoot    string
	rcServerPool    *rpcservice.RPCServerPool
	scriptEngine    *scriptEngine
	rcServicesMap   *servicesMap
}

func NewAPPServer() *appServer {
	var err error
	app := &appServer{}
	app.dataMap = utility.NewDataMap()
	app.dataMap.Set("domain", zkClient.Domain)
	app.dataMap.Set("ip", zkClient.LocalIP)
	app.Log, err = logger.New("app server", true)
	app.appServerConfig = app.dataMap.Translate(appServerConfig)
	app.rcServerRoot = app.dataMap.Translate(rcServerRoot)
	app.rcServerPool = rpcservice.NewRPCServerPool()
	app.scriptEngine = NewScriptEngine(app)
	app.rcServicesMap = NewServiceMap()
	if err != nil {
		log.Print(err)
	}
	return app
}
