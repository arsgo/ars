package main

import (
	"encoding/json"
	"runtime/debug"
	"time"

	"github.com/colinyl/ars/monitor"
)

//RCSnap RC server快照信息
type RCSnap struct {
	Domain  string                  `json:"domain"`
	Path    string                  `json:"path"`
	Address string                  `json:"address"`
	Server  string                  `json:"server"`
	Last    string                  `json:"last"`
	Sys     *monitor.SysMonitorInfo `json:"sys"`
	ip      string
}

//GetSnap 获取RC服务的快照信息
func (rs RCSnap) GetSnap() string {
	snap := rs
	snap.Last = time.Now().Format("20060102150405")
	snap.Sys, _ = monitor.GetSysMonitorInfo()
	buffer, _ := json.Marshal(&snap)
	//fmt.Println(string(buffer))
	return string(buffer)
}

//StartRefreshSnap 启动定时刷新
func (rc *RCServer) StartRefreshSnap() {
	defer rc.recover()
	tp := time.NewTicker(time.Second * 60)
	free := time.NewTicker(time.Second * 120)
	for {
		select {
		case <-tp.C:
			rc.clusterClient.ResetSnap(rc.snap.Path, rc.snap.GetSnap())
		case <-free.C:
			debug.FreeOSMemory()

		}
	}

}
