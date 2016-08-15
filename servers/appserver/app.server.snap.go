package main

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/lib4go/sysinfo"
)

type ExtSnap struct {
	Script json.RawMessage `json:"script"`
	RPC    json.RawMessage `json:"rpc"`
}

//AppSnap  app server快照信息
type AppSnap struct {
	appserver *AppServer
	Address   string               `json:"address"`
	Server    string               `json:"server"`
	Version   string               `json:"version"`
	Last      string               `json:"last"`
	Mem       uint64               `json:"mem"`
	Sys       *base.SysMonitorInfo `json:"sys"`
	Snap      ExtSnap              `json:"snap"`
	ip        string
	port      string
	path      string
}

//GetSnap 获取快照信息
func (as AppSnap) GetSnap() string {
	snap := as
	snap.Last = time.Now().Format("20060102150405")
	snap.Sys, _ = base.GetSysMonitorInfo()
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
	snap.Sys, _ = base.GetSysMonitorInfo()
	buffer, _ := json.Marshal(&snap)
	return string(buffer)
}

func (app *AppServer) recover() {
	//if r := recover(); r != nil {
	//	app.Log.Error(r, string(debug.Stack()))
	//	}
}

//StartRefreshSnap 启动定时刷新快照信息
func (app *AppServer) StartRefreshSnap() {
	defer app.recover()
	app.snap.path, _ = app.clusterClient.CreateAppServer(app.snap.port, app.snap.GetSnap())
	app.ResetJobSnap()
	tp := time.NewTicker(time.Second * 60)
	free := time.NewTicker(time.Second * 302)
	for {
		select {
		case <-tp.C:
			app.Log.Info(" -> 更新app server快照信息")
			app.ResetAPPSnap()
			app.ResetJobSnap()
		case <-free.C:
			app.Log.Infof(" -> 清理内存...%dM", sysinfo.GetAPPMemory())
			debug.FreeOSMemory()
		}
	}
}
func (app *AppServer) resetAppServer() {
	app.Log.Debug(" -> 更新所有服务")
	app.CloseAppServer()
	app.CloseJobSnap()
	time.Sleep(time.Second)
	app.ResetAPPSnap()
	app.ResetJobSnap()
}

//ResetJobSnap 重置JOB快照信息
func (app *AppServer) ResetJobSnap() (err error) {
	paths := app.scriptPorxy.GetTasks()
	for _, path := range paths {
		app.clusterClient.SetNode(path, app.snap.GetJobSnap(app.jobServer.Address))
	}
	return nil
}
func (app *AppServer) CloseJobSnap() (err error) {
	paths := app.scriptPorxy.GetTasks()
	for _, path := range paths {
		app.clusterClient.CloseNode(path)
	}
	return nil
}

//ResetAPPSnap 刷新APP快照信息
func (app *AppServer) ResetAPPSnap() (err error) {
	return app.clusterClient.SetNode(app.snap.path, app.snap.GetSnap())
}

//CloseAppServer 关闭 APP Server
func (app *AppServer) CloseAppServer() (err error) {
	return app.clusterClient.CloseAppServer(app.snap.path)
}
