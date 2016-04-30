package cluster

import (
	"sync"

	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/ars/webservice"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

const (
	appServerConfig  = "@domain/app/config/@ip"
	appServerPath    = "@domain/app/servers/@ip"
	jobConsumerPath  = "@domain/job/servers/@jobName/job_"
	jobConsumerValue = `{"ip":"@ip@jobPort","last":@now}`
)

type taskConfig struct {
	Trigger string `json:"trigger"`
	Script  string `json:"script"`
}
type taskRouteConfig struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Script string `json:"script"`
}

type serverConfig struct {
	ServerType string             `json:"type"`
	Routes     []*taskRouteConfig `json:"routes"`
}
type AppConfig struct {
	Status  string         `json:"status"`
	Tasks   []*taskConfig  `json:"tasks"`
	Jobs    []string       `json:"jobs"`
	Server  *serverConfig  `json:"server"`
	Monitor *monitorConfig `json:"monitor"`
}

type RCServerConfig struct {
	Domain string
	IP     string
	Port   string
	Server string
	Online string
}
type scriptHandler struct {
	data   *taskRouteConfig
	server *appServer
}

type appServer struct {
	dataMap           *utility.DataMap
	Last              int64
	Log               *logger.Logger
	zkClient          *zkClientObj
	appServerConfig   string
	rcServerRoot      string
	rcServerPool      *rpcservice.RPCServerPool
	scriptEngine      *scriptEngine
	rcServicesMap     *servicesMap
	jobServer         *rpcservice.RPCServer
	hasStartJobServer bool
	jobServerAdress   string
	appServerAddress  string
	lk                sync.Mutex
	jobNames          map[string]string
	apiServer         *webservice.WebService
	apiServerAddress  string
	appRoutes         []*taskRouteConfig
	scriptHandlers    map[string]*scriptHandler
	monitor           *serverMonitor
}

func NewAPPServer() *appServer {
	app := &appServer{}
	app.Log, _ = logger.New("app server", true)
	return app
}
func (app *appServer) init() (err error) {
	app.zkClient = NewZKClient()
	app.dataMap = app.zkClient.dataMap.Copy()
	app.appServerConfig = app.dataMap.Translate(appServerConfig)
	app.rcServerRoot = app.dataMap.Translate(rcServerRoot)
	app.appServerAddress = app.dataMap.Translate(appServerPath)
	app.rcServerPool = rpcservice.NewRPCServerPool()
	app.scriptEngine = NewScriptEngine(app)
	app.rcServicesMap = NewServiceMap()
	app.jobNames = make(map[string]string)
	app.scriptHandlers = make(map[string]*scriptHandler)
	app.monitor = NewMonitor(app.zkClient, app.Log,app)
	return
}

func (r *appServer) Start() (err error) {
	if err = r.init(); err != nil {
		return
	}
	r.WatchRCServerChange(func(config []*RCServerConfig, err error) {
		r.BindRCServer(config, err)
	})

	r.WatchConfigChange(func(config *AppConfig, err error) error {
		r.BindTask(config, err)
		return nil
	})

	return nil
}

func (r *appServer) Stop() error {
	defer func() {
		recover()
	}()

	r.zkClient.ZkCli.Close()
	if r.jobServer != nil {
		r.jobServer.Stop()
	}
	r.Log.Info("::app server closed")
	return nil
}
