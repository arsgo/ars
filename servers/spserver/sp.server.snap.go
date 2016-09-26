package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/arsgo/ars/snap"
	"github.com/arsgo/lib4go/sysinfo"
	"github.com/arsgo/lib4go/utility"
)

//SPSnap sp server快照信息
type SPSnap struct {
	spserver *SPServer
	Address  string      `json:"address"`
	Service  string      `json:"service"`
	Refresh  int         `json:"refresh"`
	AppMem   string      `json:"am"`
	Version  string      `json:"version"`
	CPU      string      `json:"cpu"`
	Mem      string      `json:"mem"`
	Disk     string      `json:"disk"`
	Last     string      `json:"last"`
	Snap     interface{} `json:"snap"`
}

//ResetSPSnap 重置SP server 快照
func (sp *SPServer) updateSnap(snaps map[string]interface{}) {
	services := sp.rpcServer.GetServicePath()
	sp.Log.Infof(" - >更新 sp server 快照信息:svs(%d), 内存...%dM", len(services), sysinfo.GetAPPMemory())
	sp.snap.AppMem = fmt.Sprintf("%dm", sysinfo.GetAPPMemory())
	sp.snap.CPU = sysinfo.GetAvaliabeCPU().Used
	sp.snap.Mem = sysinfo.GetAvaliabeMem().Used
	sp.snap.Disk = sysinfo.GetAvaliabeDisk().Used
	for k, v := range services {
		nsnap := utility.CloneMap(snaps)
		utility.Merge(nsnap, sp.rpcServerCollector.Customer(k).Get())
		nsnap[k] = sp.rpcServerCollector.GetByName(k)
		sp.clusterClient.SetNode(v, sp.snap.getSnap(k, nsnap))
	}
}

//StartRefreshSnap 启动快照刷新服务
func (sp *SPServer) startRefreshSnap() {
	defer sp.recover()
	snap.Bind(time.Second*time.Duration(sp.snap.Refresh), sp.updateSnap)
}

func (sp *SPServer) resetSPSnap() {
	sp.snapLogger.Info(" -> 更新所有服务")
	services := sp.rpcServer.GetServicePath()
	for _, v := range services {
		sp.clusterClient.CloseSPServer(v)
	}
	time.Sleep(time.Second)
	sp.updateSnap(snap.GetData())
}

//GetSnap 获取指定服务的快照信息
func (sn SPSnap) getSnap(service string, snaps map[string]interface{}) string {
	snap := sn
	snap.Service = service
	snap.Last = time.Now().Format("20060102150405")
	//	snap.Snap = snaps
	buffer, err := json.Marshal(&snap)
	if err != nil {
		sn.spserver.Log.Error("get snap:", err)
	}
	return string(buffer)
}

func (sn SPSnap) getDefSnap(service string) string {
	//cache := make(map[string]interface{})
	//cache["rpc"] = sn.spserver.rpcClient.GetSnap()
	sn.AppMem = fmt.Sprintf("%dm", sysinfo.GetAPPMemory())
	sn.CPU = sysinfo.GetAvaliabeCPU().Used
	sn.Mem = sysinfo.GetAvaliabeMem().Used
	sn.Disk = sysinfo.GetAvaliabeDisk().Used
	return sn.getSnap(service, snap.GetData())
}
