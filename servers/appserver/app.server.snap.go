package main

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/colinyl/ars/monitor"
	"github.com/colinyl/lib4go/sysinfo"
)

type ExtSnap struct {
	Script json.RawMessage `json:"script"`
	RPC    json.RawMessage `json:"rpc"`
}

//AppSnap  app server快照信息
type AppSnap struct {
	appserver *AppServer
	Address   string                  `json:"address"`
	Server    string                  `json:"server"`
	Last      string                  `json:"last"`
	Mem       uint64                  `json:"mem"`
	Sys       *monitor.SysMonitorInfo `json:"sys"`
	//	ServerSnap json.RawMessage         `json:"serverSnap"`
	Snap ExtSnap `json:"snap"`
	ip   string
}

//GetSnap 获取快照信息
func (as AppSnap) GetSnap() string {
	snap := as
	snap.Last = time.Now().Format("20060102150405")

	snap.Sys, _ = monitor.GetSysMonitorInfo()
	//if as.appserver.httpServer != nil {
	//	snap.ServerSnap, _ = json.Marshal(as.appserver.httpServer.GetSnap())
	//	} else {
	//	snap.ServerSnap = []byte("{}")
	//	}
	snap.Snap.RPC, _ = json.Marshal(as.appserver.rpcClient.GetSnap())
	snap.Snap.Script, _ = json.Marshal(as.appserver.scriptPool.GetSnap())
	snap.Mem = sysinfo.GetAPPMemory()
	buffer, _ := json.Marshal(&snap)
	r := string(buffer)
	return r
}

//GetJobSnap 获取快照信息
func (as AppSnap) GetJobSnap(server string) string {
	snap := as
	snap.Server = fmt.Sprint(snap.ip, server)
	snap.Last = time.Now().Format("20060102150405")
	snap.Sys, _ = monitor.GetSysMonitorInfo()
	buffer, _ := json.Marshal(&snap)
	return string(buffer)
}

func (app *AppServer) recover() {
	if r := recover(); r != nil {
		app.Log.Error(r, string(debug.Stack()))
	}
}

//StartRefreshSnap 启动定时刷新快照信息
func (app *AppServer) StartRefreshSnap() {
	defer app.recover()
	app.Log.Info("更新app server快照信息")
	app.ResetAPPSnap()
	app.ResetJobSnap()
	tp := time.NewTicker(time.Second * 60)
	free := time.NewTicker(time.Second * 122)
	for {
		select {
		case <-tp.C:
			app.Log.Info("更新app server快照信息")
			app.ResetAPPSnap()
			app.ResetJobSnap()
		case <-free.C:
			app.Log.Infof("清理内存...%dM", sysinfo.GetAPPMemory())
			debug.FreeOSMemory()
		}
	}
}

//ResetJobSnap 重置JOB快照信息
func (app *AppServer) ResetJobSnap() (err error) {
	paths := app.jobConsumerScriptHandler.GetTasks()
	for _, path := range paths {
		app.clusterClient.UpdateJobConsumerPath(path, app.snap.GetJobSnap(app.jobConsumerRPCServer.Address))
	}
	return nil
}

//ResetAPPSnap 刷新APP快照信息
func (app *AppServer) ResetAPPSnap() (err error) {
	snap := app.snap.GetSnap()
	err = app.clusterClient.ResetAppServerSnap(snap)
	return
}
