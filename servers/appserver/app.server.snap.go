package main

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/colinyl/ars/monitor"
)

//AppSnap  app server快照信息
type AppSnap struct {
	appserver  *AppServer
	Address    string                  `json:"address"`
	Server     string                  `json:"server"`
	Last       string                  `json:"last"`
	Sys        *monitor.SysMonitorInfo `json:"sys"`
	ScriptSnap json.RawMessage         `json:"scriptSnap"`
	RPCSnap    json.RawMessage         `json:"rpcSnap"`
	ip         string
}

//GetSnap 获取快照信息
func (as AppSnap) GetSnap() string {
	snap := as
	snap.Last = time.Now().Format("20060102150405")
	snap.Sys, _ = monitor.GetSysMonitorInfo()
	snap.RPCSnap, _ = json.Marshal(as.appserver.rpcClient.GetSnap().Snaps)
	snap.ScriptSnap, _ = json.Marshal(as.appserver.scriptPool.GetSnap().Snaps)
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
		app.Log.Fatal(r)
	}
}

//StartRefreshSnap 启动定时刷新快照信息
func (app *AppServer) StartRefreshSnap() {
	defer app.recover()
	tp := time.NewTicker(time.Second * 60)
	free := time.NewTicker(time.Second * 120)
	for {
		select {
		case <-tp.C:
			app.Log.Info("更新app server快照信息")
			app.ResetAPPSnap()
			app.ResetJobSnap()
		case <-free.C:
			app.Log.Info("清理内存...")
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
	err = app.clusterClient.ResetAppServerSnap(app.snap.GetSnap())
	return
}
