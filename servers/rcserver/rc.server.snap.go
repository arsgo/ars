package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/arsgo/ars/snap"
	"github.com/arsgo/lib4go/sysinfo"
)

//RCSnap RC server快照信息
type RCSnap struct {
	rcServer *RCServer
	Domain   string `json:"domain"`
	path     string
	Address  string      `json:"address"`
	Server   string      `json:"server"`
	Refresh  int         `json:"refresh"`
	AppMem   string      `json:"am"`
	Version  string      `json:"version"`
	CPU      string      `json:"cpu"`
	Mem      string      `json:"mem"`
	Disk     string      `json:"disk"`
	Last     string      `json:"last"`
	Snap     interface{} `json:"snap"`
	Cache    interface{} `json:"cache"`
}

//GetServicesSnap 获取RC服务的快照信息
func (rs RCSnap) GetServicesSnap(services map[string]interface{}) string {
	snap := rs
	snap.Last = time.Now().Format("20060102150405")
	snap.AppMem = fmt.Sprintf("%dm", sysinfo.GetAPPMemory())
	snap.CPU = sysinfo.GetAvaliabeCPU().Used
	snap.Mem = sysinfo.GetAvaliabeMem().Used
	snap.Disk = sysinfo.GetAvaliabeDisk().Used

	rpcs := rs.rcServer.rpcServerCollector.Get()
	if len(rpcs) > 0 {
		services["rpc"] = rpcs
	}
	schedulers := rs.rcServer.schedulerCollector.Get()
	if len(schedulers) > 0 {
		services["jobs"] = schedulers
	}
	cache := make(map[string]interface{})
	cache["rpc"] = rs.rcServer.spRPCClient.GetSnap()
	snap.Cache = cache
	snap.Snap = services

	buffer, err := json.Marshal(&snap)
	if err != nil {
		rs.rcServer.Log.Error("更新快照异常：", err)
	}
	return string(buffer)
}

//startRefreshSnap 启动定时刷新
func (rc *RCServer) startRefreshSnap() {
	defer rc.recover()
	snap.Bind(time.Second*time.Duration(rc.snap.Refresh), rc.updateSnap)
}

func (rc *RCServer) setDefSnap() {
	rc.updateSnap(snap.GetData())
}

func (rc *RCServer) updateSnap(services map[string]interface{}) {
	rc.snapLogger.Info(" -> 更新 rc server快照信息")
	rc.clusterClient.SetNode(rc.snap.path, rc.snap.GetServicesSnap(services))
}
