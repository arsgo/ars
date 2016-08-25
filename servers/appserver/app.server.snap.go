package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/arsgo/ars/snap"
	"github.com/arsgo/lib4go/sysinfo"
)

type ExtSnap struct {
	Script json.RawMessage `json:"script"`
	RPC    json.RawMessage `json:"rpc"`
}

//AppSnap  app server快照信息
type AppSnap struct {
	appserver *AppServer
	Address   string      `json:"address"`
	Server    string      `json:"server"`
	Refresh   int         `json:"refresh"`
	Version   string      `json:"version"`
	CPU       interface{} `json:"cpu"`
	Mem       interface{} `json:"mem"`
	Disk      interface{} `json:"disk"`
	Last      string      `json:"last"`
	Snap      interface{} `json:"snap"`
	Cache     interface{} `json:"cache"`

	port string
	path string
}

type MQSnap struct {
	Address string      `json:"address"`
	Name    string      `json:"name"`
	Refresh int         `json:"refresh"`
	Version string      `json:"version"`
	CPU     interface{} `json:"cpu"`
	Mem     interface{} `json:"mem"`
	Disk    interface{} `json:"disk"`
	Last    string      `json:"last"`
	Snap    interface{} `json:"snap"`
}

type JobConsumerSnap struct {
	Address string      `json:"address"`
	Server  string      `json:"server"`
	Refresh int         `json:"refresh"`
	Version string      `json:"version"`
	CPU     interface{} `json:"cpu"`
	Mem     interface{} `json:"mem"`
	Disk    interface{} `json:"disk"`
	Last    string      `json:"last"`
	Snap    interface{} `json:"snap"`
}

type JobLocalSnap struct {
	Address string      `json:"address"`
	Refresh int         `json:"refresh"`
	Version string      `json:"version"`
	CPU     interface{} `json:"cpu"`
	Mem     interface{} `json:"mem"`
	Disk    interface{} `json:"disk"`
	Last    string      `json:"last"`
	Snap    interface{} `json:"snap"`
}

//GetSnap 获取快照信息
func (as AppSnap) GetSnap() string {
	snap := as
	snap.Last = time.Now().Format("20060102150405")
	buffer, _ := json.Marshal(&snap)
	r := string(buffer)
	return r
}

//getJobConsumerSnap 获取快照信息
func (as AppSnap) getJobConsumerSnap(server string) string {
	var snap JobConsumerSnap
	snap.Address = as.Address
	snap.Refresh = as.Refresh
	snap.Version = as.Version
	snap.CPU = as.CPU
	snap.Mem = as.Mem
	snap.Disk = as.Disk
	snap.Snap = as.Snap
	snap.Server = fmt.Sprint(as.appserver.conf.IP, server)
	snap.Last = time.Now().Format("20060102150405")
	buffer, _ := json.Marshal(&snap)
	return string(buffer)
}

//getJobLocalSnap 获取快照信息
func (as AppSnap) getJobLocalSnap() string {
	var snap JobLocalSnap
	snap.Address = as.Address
	snap.Refresh = as.Refresh
	snap.Version = as.Version
	snap.CPU = as.CPU
	snap.Mem = as.Mem
	snap.Disk = as.Disk
	snap.Snap = as.Snap
	snap.Last = time.Now().Format("20060102150405")
	buffer, _ := json.Marshal(&snap)
	return string(buffer)
}

//geMQSnap 获取快照信息
func (as AppSnap) getMQSnap(name string) string {
	var snap MQSnap
	snap.Name = name
	snap.Address = as.Address
	snap.Refresh = as.Refresh
	snap.Version = as.Version
	snap.CPU = as.CPU
	snap.Mem = as.Mem
	snap.Disk = as.Disk
	snap.Snap = as.Snap
	snap.Last = time.Now().Format("20060102150405")
	buffer, _ := json.Marshal(&snap)
	return string(buffer)
}

//startRefreshSnap 启动定时刷新快照信息
func (app *AppServer) startRefreshSnap() {
	defer app.recover()
	app.updateSnap(snap.GetData())
	snap.Bind(time.Second*time.Duration(app.snap.Refresh), app.updateSnap)
}
func (app *AppServer) resetAppServer() {
	app.closeSnap()
	time.Sleep(time.Second)
	app.updateSnap(snap.GetData())
}

//updateSnap 重置JOB快照信息
func (app *AppServer) updateSnap(services map[string]interface{}) {
	app.snapLogger.Info(" -> 更新 app server快照信息")
	cache := make(map[string]interface{})
	app.appendSystemSnap(services, cache)
	app.snap.Snap = services
	if app.jobServer.Available {
		jobPaths := app.jobServer.GetServicePath()
		for name, path := range jobPaths {
			services["job.consumer"] = app.jobServerCollector.GetByName(name)
			app.clusterClient.SetNode(path, app.snap.getJobConsumerSnap(app.jobServer.Address))
		}
	}
	delete(services, "job.consumer")
	if app.mqService.Available {
		mqPaths := app.mqService.GetServices()
		for name, path := range mqPaths {
			services["mq.consumer"] = app.mqConsumerCollector.GetByName(name)
			app.clusterClient.SetNode(path, app.snap.getMQSnap(name))
		}
	}

	delete(services, "mq.consumer")
	locajobPath := app.localJobPaths.GetAll()
	for name, p := range locajobPath {
		services["job.local"] = app.jobLocalCollector.GetByName(name)
		app.clusterClient.SetNode(p.(string), app.snap.getJobLocalSnap())
	}
	delete(services, "job.local")
	if app.apiServer != nil && app.apiServer.Available {
		services["api.server"] = app.apiServerCollector.Get()
		services["customer"] = app.scriptCollector.Get()
		cache["rpc"] = app.rpcClient.GetSnap()
		if strings.EqualFold(app.snap.path, "") {
			app.snap.path, _ = app.clusterClient.CreateAppServer(app.snap.port, app.snap.GetSnap())
		} else {
			app.clusterClient.SetNode(app.snap.path, app.snap.GetSnap())
		}
	}
}

func (app *AppServer) closeSnap() {
	app.clusterClient.CloseAppServer(app.snap.path)
	paths := app.jobServer.GetServicePath()
	for _, path := range paths {
		app.clusterClient.CloseNode(path)
	}
	mqPaths := app.mqService.GetServices()
	for _, path := range mqPaths {
		app.clusterClient.CloseNode(path)
	}
	if !strings.EqualFold(app.snap.path, "") {
		app.clusterClient.CloseNode(app.snap.path)
		app.snap.path = ""
	}
}
func (app *AppServer) appendSystemSnap(snaps map[string]interface{}, cache map[string]interface{}) {
	app.snap.CPU = sysinfo.GetAvaliabeCPU().Used
	app.snap.Mem = sysinfo.GetAvaliabeMem().Used
	app.snap.Disk = sysinfo.GetAvaliabeDisk().Used
	app.snap.Snap = snaps
	app.snap.Cache = cache
}
