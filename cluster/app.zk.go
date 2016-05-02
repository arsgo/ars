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

func (d *appServer) CreateJobSnap(jobMap map[string]string) {
	d.jobPaths.mutex.Lock()
	defer d.jobPaths.mutex.Unlock()
	dmap := d.dataMap.Copy()
	for i := range jobMap {
		if _, ok := d.jobPaths.data[i]; !ok {
			dmap.Set("jobName", i)
			path, err := d.zkClient.ZkCli.CreateSeqNode(dmap.Translate(jobConsumerPath),
				d.snap.GetSnap())
			if err != nil {
				d.Log.Error(err)
				continue
			}
			d.jobPaths.data[i] = path
			d.Log.Infof("::start job service:%s", i)
		}
	}
}

func (d *appServer) ResetJobSnap() (err error) {
	paths := d.jobPaths.getData()
	for _, path := range paths {
		d.zkClient.ZkCli.UpdateValue(path, d.snap.GetSnap())
	}
	return nil
}
func (d *appServer) ResetAPPSnap() (err error) {
	if d.zkClient.ZkCli.Exists(d.appServerAddress) {
		err = d.zkClient.ZkCli.UpdateValue(d.appServerAddress, d.snap.GetSnap())
	} else {
		err = d.zkClient.ZkCli.CreatePath(d.appServerAddress, d.snap.GetSnap())
	}
	return
}
