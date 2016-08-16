package main

import (
	"encoding/json"
	"runtime/debug"
	"sync"
	"time"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/lib4go/sysinfo"
)

type ExtSnap struct {
	Script json.RawMessage `json:"script"`
	RPC    json.RawMessage `json:"rpc"`
}

//SPSnap sp server快照信息
type SPSnap struct {
	Address string               `json:"address"`
	Service string               `json:"service"`
	Version string               `json:"version"`
	Last    string               `json:"last"`
	Mem     uint64               `json:"mem"`
	Sys     *base.SysMonitorInfo `json:"sys"`
	//ServerSnap json.RawMessage         `json:"serverSnap"`
	Snap  ExtSnap `json:"snap"`
	ip    string
	mutex sync.Mutex
}

//GetSnap 获取指定服务的快照信息
func (sn SPSnap) GetSnap(service string) string {
	sn.mutex.Lock()
	snap := sn
	sn.mutex.Unlock()
	snap.Service = service
	snap.Last = time.Now().Format("20060102150405")
	snap.Sys, _ = base.GetSysMonitorInfo()
	buffer, _ := json.Marshal(&snap)
	return string(buffer)
}

//ResetSPSnap 重置SP server 快照
func (sp *SPServer) ResetSPSnap() {
	services := sp.rpcScriptProxy.GetTasks()
	//sp.snap.ServerSnap, _ = json.Marshal(sp.rpcServer.GetSnap())
	sp.snap.Snap.RPC, _ = json.Marshal(sp.rpcClient.GetSnap())
	sp.snap.Snap.Script, _ = json.Marshal(sp.scriptPool.GetSnap())
	sp.snap.Mem = sysinfo.GetAPPMemory()
	for k, v := range services {
		sp.clusterClient.SetNode(v, sp.snap.GetSnap(k))
	}
}

func (sp *SPServer) CloseSPServer() {
	services := sp.rpcScriptProxy.GetTasks()
	for _, v := range services {
		sp.clusterClient.CloseSPServer(v)
	}
}

//StartRefreshSnap 启动快照刷新服务
func (sp *SPServer) startRefreshSnap() {
	defer sp.recover()
	sp.ResetSPSnap()
	tp := time.NewTicker(time.Second * 60)
	free := time.NewTicker(time.Second * 122)
	defer tp.Stop()
	for {
		select {
		case <-tp.C:
			sp.snapLogger.Info("更新sp server快照信息")
			sp.ResetSPSnap()
		case <-free.C:
			sp.snapLogger.Infof("清理内存...%dM", sysinfo.GetAPPMemory())
			debug.FreeOSMemory()
		}
	}
}

func (sp *SPServer) resetCluster() {
	sp.snapLogger.Info(" -> 更新所有服务")
	sp.CloseSPServer()
	time.Sleep(time.Second)
	sp.ResetSPSnap()

}
func (sp *SPServer) recover() {
	//	if r := recover(); r != nil {
	///	sp.Log.Fatal(r, string(debug.Stack()))
	//	}
}
