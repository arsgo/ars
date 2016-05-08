package main

import (
	"encoding/json"
	"time"

	"github.com/colinyl/ars/monitor"
)

//SPSnap sp server快照信息
type SPSnap struct {
	Address string                  `json:"address"`
	Service string                  `json:"service"`
	Last    string                  `json:"last"`
	Sys     *monitor.SysMonitorInfo `json:"sys"`
}

//GetSnap 获取指定服务的快照信息
func (sn SPSnap) GetSnap(service string) string {
	snap := sn
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
	tp := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-tp.C:
			sp.ResetSPSnap()
		}
	}

}