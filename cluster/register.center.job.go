package cluster

import (
	"encoding/json"
	"time"
)

//WatchJobChange 监控JOB配置变化
func (d *rcServer) WatchJobChange(callback func(config *JobConfigs, err error)) {
	if !d.IsMasterServer {
		return
	}
	d.jobCallback = callback
	d.zkClient.waitZKPathExists(d.jobConfigPath, time.Hour*8640, func(exists bool) {
		if exists {
			callback(d.getJobConfigs(d.jobConfigPath))
		} else {
			d.Log.Info("job config path not exists")
		}
	})
	d.Log.Info("::watch job config changes")
	d.zkClient.watchZKValueChange(d.jobConfigPath, func() {
		d.Log.Info("job config has changed")
		callback(d.getJobConfigs(d.jobConfigPath))
	})
}
func (d *rcServer) getJobConsumers(jobName string) (jobIPCollection []string) {
	dmap := d.dataMap.Copy()
	dmap.Set("jobName", jobName)
	path := dmap.Translate(jobConsumerRoot)
	children, err := d.zkClient.ZkCli.GetChildren(path)
	if err != nil {
		d.Log.Error(err)
		return
	}

	for _, v := range children {
		dmap.Set("path", v)
		values, err := d.zkClient.ZkCli.GetValue(dmap.Translate(jobConsumerPath))
		if err != nil {
			d.Log.Error(err)
			continue
		}
		consumer := &JobConsumerValue{}
		err = json.Unmarshal([]byte(values), &consumer)
		if err != nil {
			d.Log.Error(err)
			continue
		}
		jobIPCollection = append(jobIPCollection, consumer.IP)
	}
	return

}

func (d *rcServer) getJobConfigs(path string) (defConfigs *JobConfigs, err error) {
	defConfigs = &JobConfigs{}
	defConfigs.Jobs = make(map[string]JobConfigItem)
	if !d.zkClient.ZkCli.Exists(path) {
		return
	}
	value, err := d.zkClient.ZkCli.GetValue(path)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(value), &defConfigs.Jobs)
	return
}
