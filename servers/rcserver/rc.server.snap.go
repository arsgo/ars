package main

import (
	"encoding/json"
	"runtime/debug"
	"time"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/lib4go/sysinfo"
)

type ExtSnap struct {
	Server json.RawMessage `json:"server"`
	RPC    json.RawMessage `json:"rpc"`
}

//RCSnap RC server快照信息
type RCSnap struct {
	rcServer *RCServer
	Domain   string               `json:"domain"`
	Path     string               `json:"path"`
	Address  string               `json:"address"`
	Server   string               `json:"server"`
	Version  string               `json:"version"`
	Last     string               `json:"last"`
	Mem      uint64               `json:"mem"`
	Sys      *base.SysMonitorInfo `json:"sys"`
	Snap     ExtSnap              `json:"snap"`
	ip       string
}

//GetSnap 获取RC服务的快照信息
func (rs RCSnap) GetSnap() string {
	snap := rs
	snap.Last = time.Now().Format("20060102150405")
	snap.Sys, _ = base.GetSysMonitorInfo()
	snap.Snap.Server, _ = json.Marshal(rs.rcServer.rcRPCServer.GetSnap())
	snap.Snap.RPC, _ = json.Marshal(rs.rcServer.spRPCClient.GetSnap())
	snap.Mem = sysinfo.GetAPPMemory()
	buffer, _ := json.Marshal(&snap)
	return string(buffer)
}

//RefreshSnap 刷新快照信息
func (rc *RCServer) RefreshSnap() {
	rc.clusterClient.SetNode(rc.snap.Path, rc.snap.GetSnap())
}

//startRefreshSnap 启动定时刷新
func (rc *RCServer) startRefreshSnap() {
	defer rc.recover()
	rc.clusterClient.SetNode(rc.snap.Path, rc.snap.GetSnap())
	rc.RefreshSnap()
	tp := time.NewTicker(time.Second * 60)
	free := time.NewTicker(time.Second * 122)
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
