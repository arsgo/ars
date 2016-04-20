package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

func (d *appServer) WatchConfigChange(callback func(config *AppConfig, err error) error) {
	d.zkClient.waitZKPathExists(d.appServerConfig, time.Hour*8640, func(exists bool) {
		if !exists {
			d.Log.Infof("app config not exists:%s", d.appServerConfig)
		} else {
			callback(d.zkClient.getAppConfig(d.appServerConfig))
		}
	})
	d.Log.Info("::watch for app config changes")
	d.zkClient.watchZKValueChange(d.appServerConfig, func() {
		d.Log.Info("app config has changed")
		callback(d.zkClient.getAppConfig(d.appServerConfig))
	})
}

func (d *appServer) WatchRCServerChange(callback func([]*RCServerConfig, error)) {
	d.zkClient.waitZKPathExists(d.rcServerRoot, time.Hour*8640, func(exists bool) {
		if !exists {
			d.Log.Infof("rc server not exists:%s", d.appServerConfig)
		} else {
			callback(d.zkClient.getRCServer(d.dataMap))
		}
	})
	d.Log.Info("::watch for rc server changes")
	d.zkClient.watchZKChildrenPathChange(d.rcServerRoot, func() {
		d.Log.Info("rc server has changed")
		callback(d.zkClient.getRCServer(d.dataMap))
	})
}

func (d *appServer) resetSnap(path string, config *AppConfig) (err error) {
	data, err := json.Marshal(config)
	if err != nil {
		return
	}
	dst := new(bytes.Buffer)
	json.Indent(dst, data, "  ", "  ")
	content := fmt.Sprintf("%s", dst)
	if d.zkClient.ZkCli.Exists(path) {
		err = d.zkClient.ZkCli.UpdateValue(path, content)
	} else {
		err = d.zkClient.ZkCli.CreatePath(path, content)
	}
	return
}
