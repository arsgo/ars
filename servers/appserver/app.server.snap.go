package appserver

import (
	"encoding/json"
	"time"

	"github.com/colinyl/ars/monitor"
)

//AppSnap  app server快照信息
type AppSnap struct {
	Address string                  `json:"address"`
	Last    string                  `json:"last"`
	Sys     *monitor.SysMonitorInfo `json:"sys"`
}

//GetSnap 获取快照信息
func (as AppSnap) GetSnap() string {
	snap := as
	snap.Last = time.Now().Format("20060102150405")
	snap.Sys, _ = monitor.GetSysMonitorInfo()
	buffer, _ := json.Marshal(&snap)
	return string(buffer)
}

//StartRefreshSnap 启动定时刷新快照信息
func (app *AppServer) StartRefreshSnap() {
	tp := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-tp.C:
			app.ResetAPPSnap()
			app.ResetJobSnap()
		}
	}
}

//ResetJobSnap 重置JOB快照信息
func (app *AppServer) ResetJobSnap() (err error) {
	paths := app.jobConsumerScriptHandler.GetTasks()
	for _, path := range paths {
		app.clusterClient.UpdateJobConsumerPath(path, app.snap.GetSnap())
	}
	return nil
}

//ResetAPPSnap 刷新APP快照信息
func (app *AppServer) ResetAPPSnap() (err error) {
	err = app.clusterClient.ResetAppServerSnap(app.snap.GetSnap())
	return
}
