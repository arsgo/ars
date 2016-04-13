package cluster

import "time"

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
