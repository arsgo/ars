package cluster

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/colinyl/ars/config"
	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
	"github.com/colinyl/lib4go/webserver"
)

const (
	appServerConfig = "@domain/app/config/@ip"
	appServerPath   = "@domain/app/servers/@ip"
	jobConsumerPath = "@domain/job/servers/@jobName/job_"

	//jobConsumerValue = `{"ip":"@ip@jobPort","last":@now}`
)

type appSnap struct {
	Address string          `json:"address"`
	Last    string          `json:"last"`
	Sys     *sysMonitorInfo `json:"sys"`
}

func (a appSnap) GetSnap() string {
	snap := a
	snap.Last = time.Now().Format("20060102150405")
	snap.Sys, _ = GetSysMonitorInfo()
	buffer, _ := json.Marshal(&snap)
	return string(buffer)
}

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
	Status string        `json:"status"`
	Tasks  []*taskConfig `json:"tasks"`
	Jobs   []string      `json:"jobs"`
	Server *serverConfig `json:"server"`
}

type RCServerConfig struct {
	Domain  string
	Address string
	Server  string
}
type scriptHandler struct {
	data   *taskRouteConfig
	server *appServer
}
type jobPaths struct {
	data  map[string]string
	mutex sync.Mutex
}

func (j jobPaths) getData() map[string]string {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	d := make(map[string]string)
	return utility.MergeStringMap(j.data, d)
}

type appServer struct {
	dataMap utility.DataMap
	//Last              int64
	Log               *logger.Logger
	zkClient          *clusterClient
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
	jobPaths          jobPaths
	apiServer         *webserver.WebServer
	apiServerAddress  string
	appRoutes         []*taskRouteConfig
	scriptHandlers    map[string]*scriptHandler
	snap              appSnap
}

func NewAPPServer() *appServer {
	app := &appServer{}
	app.Log, _ = logger.New("app server", true)
	return app
}
func (app *appServer) init() (err error) {
	app.zkClient = NewClusterClient()
	app.dataMap = app.zkClient.dataMap.Copy()
	app.appServerConfig = app.dataMap.Translate(appServerConfig)
	app.rcServerRoot = app.dataMap.Translate(rcServerRoot)
	app.appServerAddress = app.dataMap.Translate(appServerPath)
	app.rcServerPool = rpcservice.NewRPCServerPool()
	app.scriptEngine = NewScriptEngine(app)
	app.rcServicesMap = NewServiceMap()
	app.scriptHandlers = make(map[string]*scriptHandler)
	app.snap = appSnap{Address: config.Get().IP}
	app.jobPaths = jobPaths{data: make(map[string]string)}
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
	go r.StartRefreshSnap()
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
