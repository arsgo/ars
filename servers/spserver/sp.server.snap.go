package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/arsgo/ars/snap"
	"github.com/arsgo/lib4go/sysinfo"
)


//SPSnap sp server快照信息
type SPSnap struct {
	spserver *SPServer
	Address  string      `json:"address"`
	Service  string      `json:"service"`
	Version  string      `json:"version"`
	CPU      string      `json:"cpu"`
	Mem      string      `json:"mem"`
	Disk     string      `json:"disk"`
	Last     string      `json:"last"`
	Snap     interface{} `json:"snap"`
	Cache    interface{} `json:"cache"`
}

//ResetSPSnap 重置SP server 快照
func (sp *SPServer) ResetSPSnap(snaps map[string]interface{}) {
	sp.snapLogger.Info("更新sp server快照信息")
	services := sp.rpcScriptProxy.GetTasks()
	cache := make(map[string]interface{})
	cache["rpc"] = sp.rpcClient.GetSnap()
	for k, v := range services {
		snaps["rpc.server"] = sp.rpcServerCollector.GetByName(v)
		sp.clusterClient.SetNode(v, sp.snap.getSnap(k, snaps, cache))
	}
}

//StartRefreshSnap 启动快照刷新服务
func (sp *SPServer) startRefreshSnap() {
	defer sp.recover()
	snap.Bind(time.Second*60, sp.ResetSPSnap)
}

func (sp *SPServer) resetSPSnap() {
	sp.snapLogger.Info(" -> 更新所有服务")
	services := sp.rpcScriptProxy.GetTasks()
	for _, v := range services {
		sp.clusterClient.CloseSPServer(v)
	}
	time.Sleep(time.Second)
	sp.ResetSPSnap(make(map[string]interface{}))
}


//GetSnap 获取指定服务的快照信息
func (sn SPSnap) getSnap(service string, snaps map[string]interface{}, cache map[string]interface{}) string {
	snap := sn
	snap.Service = service
	snap.Last = time.Now().Format("20060102150405")
	snap.Snap = snaps
	snap.Cache = cache
	snap.CPU = sysinfo.GetAvaliabeCPU().Used
	snap.Mem = sysinfo.GetAvaliabeMem().Used
	snap.Disk = sysinfo.GetAvaliabeDisk().Used
	buffer, err := json.Marshal(&snap)
	if err != nil {
		fmt.Println(err)
	}
	return string(buffer)
}

func (sn SPSnap) getDefSnap(service string) string {
	snaps := make(map[string]interface{})
	cache := make(map[string]interface{})
	cache["rpc"] = sn.spserver.rpcClient.GetSnap()
	return sn.getSnap(service, snaps, cache)
}
