package main

import (
	"encoding/json"
	"runtime/debug"
	"time"

	"github.com/colinyl/ars/monitor"
	"github.com/colinyl/lib4go/sysinfo"
)

//RCSnap RC server快照信息
type RCSnap struct {
	rcServer *RCServer
	Domain   string                  `json:"domain"`
	Path     string                  `json:"path"`
	Address  string                  `json:"address"`
	Server   string                  `json:"server"`
	Last     string                  `json:"last"`
	Mem      uint64                  `json:"mem"`
	Sys      *monitor.SysMonitorInfo `json:"sys"`
	RPCSnap  json.RawMessage         `json:"rpcSnap"`
	ip       string
}

//GetSnap 获取RC服务的快照信息
func (rs RCSnap) GetSnap() string {
	snap := rs
	snap.Last = time.Now().Format("20060102150405")
	snap.Sys, _ = monitor.GetSysMonitorInfo()
	snap.RPCSnap, _ = json.Marshal(rs.rcServer.spRPCClient.GetSnap().Snaps)
	snap.Mem = sysinfo.GetAPPMemory()
	buffer, _ := json.Marshal(&snap)
	return string(buffer)
}

//RefreshSnap 刷新快照信息
func (rc *RCServer) RefreshSnap() {
	rc.clusterClient.ResetSnap(rc.snap.Path, rc.snap.GetSnap())
}

//StartRefreshSnap 启动定时刷新
func (rc *RCServer) StartRefreshSnap() {
	defer rc.recover()
	rc.clusterClient.ResetSnap(rc.snap.Path, rc.snap.GetSnap())
	tp := time.NewTicker(time.Second * 60)
	free := time.NewTicker(time.Second * 120)
	for {
		select {
		case <-tp.C:
			rc.Log.Info("更新RC Server快照信息")
			rc.RefreshSnap()
		case <-free.C:
			rc.Log.Infof("清理内存...%dM", sysinfo.GetAPPMemory())
			debug.FreeOSMemory()

		}
	}

}
