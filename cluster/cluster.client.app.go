package cluster

import "encoding/json"

//WatchAppTaskChange 监控APP Server的配置文件变化
func (client *ClusterClient) WatchAppTaskChange(callback func(config *AppServerStartupConfig, err error) error) {
	client.WaitClusterPathExists(client.appServerTaskPath, client.timeout, func(exists bool) {
		if !exists {
			client.Log.Infof("app config not exists:%s", client.appServerTaskPath)
		} else {
			callback(client.GetAppServerStartupConfig(client.appServerTaskPath))
		}
	})
	client.Log.Info("::watch for app config changes")
	client.WatchClusterValueChange(client.appServerTaskPath, func() {
		client.Log.Info("app config has changed")
		callback(client.GetAppServerStartupConfig(client.appServerTaskPath))
	})
}

//GetAppServerStartupConfig 获取App Server的配置数据
func (client *ClusterClient) GetAppServerStartupConfig(path string) (config *AppServerStartupConfig, err error) {
	config = &AppServerStartupConfig{}
	values, err := client.handler.GetValue(path)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(values), &config)
	return
}
