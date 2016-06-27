package main

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/colinyl/ars/monitor"
)

//SPSnap sp server快照信息
type SPSnap struct {
	Address string                  `json:"address"`
	Service string                  `json:"service"`
	Last    string                  `json:"last"`
	Sys     *monitor.SysMonitorInfo `json:"sys"`
	ip      string
	mutex   sync.Mutex
}

//GetSnap 获取指定服务的快照信息
func (sn SPSnap) GetSnap(service string) string {
	sn.mutex.Lock()
	snap := sn
	sn.mutex.Unlock()
	snap.Service = service
	snap.Last = time.Now().Format("20060102150405")
	snap.Sys, _ = monitor.GetSysMonitorInfo()
	buffer, _ := json.Marshal(&snap)
	return string(buffer)
}

//ResetSPSnap 重置SP server 快照
func (sp *SPServer) ResetSPSnap() {
	services := sp.rpcScriptProxy.GetTasks()
	for k, v := range services {
		sp.clusterClient.ResetSnap(v, sp.snap.GetSnap(k))
	}
}

//StartRefreshSnap 启动快照刷新服务
func (sp *SPServer) StartRefreshSnap() {
	defer sp.recover()
	tp := time.NewTicker(time.Second * 60)
	defer tp.Stop()
	for {
		select {
		case <-tp.C:
			sp.ResetSPSnap()
		}
	}

}
func (sp *SPServer) recover() {
	if r := recover(); r != nil {
		sp.Log.Fatal(r)
	}
}
