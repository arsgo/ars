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
	waitZKPathExists(d.jobConfigPath, time.Hour*8640, func(exists bool) {
		if exists {
			callback(getJobConfigs(d.jobConfigPath))
		} else {
			d.Log.Info("job config path not exists")
		}
	})
	d.Log.Info("::watch job config changes")
	watchZKValueChange(d.jobConfigPath, func() {
        d.Log.Info("job config has changed")
		callback(getJobConfigs(d.jobConfigPath))
	})
}

func getJobConfigs(path string) (defConfigs *JobConfigs, err error) {
	defConfigs = &JobConfigs{}
	defConfigs.Jobs = make(map[string]JobConfigItem)
	if !zkClient.ZkCli.Exists(path) {
		return
	}
	value, err := zkClient.ZkCli.GetValue(path)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(value), &defConfigs.Jobs)
	return
}
