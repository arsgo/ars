package cluster

import "encoding/json"
import "runtime/debug"

//recover 从异常中恢复
func (client *ClusterClient) recover() {
	if r := recover(); r != nil {
		client.Log.Fatal(r, string(debug.Stack()))
	}
}

//WatchAppTaskChange 监控APP Server的配置文件变化
func (client *ClusterClient) WatchAppTaskChange(callback func(config *AppServerStartupConfig, err error) error) {
	client.WaitClusterPathExists(client.appServerTaskPath, client.timeout, func(exists bool) {
		if !exists {
			client.Log.Infof("app config not exists:%s", client.appServerTaskPath)
		} else {
			go func() {
				defer client.recover()
				callback(client.GetCurrentAppServerTask())
			}()
		}
	})
	client.Log.Info("::watch for app config changes")
	client.WatchClusterValueChange(client.appServerTaskPath, func() {
		client.Log.Info(" -> app config has changed")
		go func() {
			defer client.recover()
			callback(client.GetCurrentAppServerTask())
		}()
	})
}

//UpdateAppServerTask 更新AppServer配置文件
func (client *ClusterClient) UpdateAppServerTask(ip string, config *AppServerStartupConfig) (err error) {
	nmap := client.dataMap.Copy()
	nmap.Set("ip", ip)
	path := nmap.Translate(p_appTaskConfig)

	buffer, err := json.Marshal(config)
	if err != nil {
		return
	}
	err = client.handler.UpdateValue(path, string(buffer))
	return
}

//GetCurrentAppServerTask 获取当前AppServer任务
func (client *ClusterClient)GetCurrentAppServerTask() (config *AppServerStartupConfig, err error){
	return client.GetAppServerTask(client.IP)
}
//GetAppServerTask 获取App Server的配置数据
func (client *ClusterClient) GetAppServerTask(ip string) (config *AppServerStartupConfig, err error) {
	nmap := client.dataMap.Copy()
	nmap.Set("ip", ip)
	path := nmap.Translate(p_appTaskConfig)

	config = &AppServerStartupConfig{}
	values, err := client.handler.GetValue(path)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(values), &config)
	if err != nil {
		return
	}
	/*jobs, err := client.GetJobTask()
	if err != nil {
		return
	}
	for _, v := range config.JobNames {
		if _, ok := jobs[v]; ok {
			config.Jobs = append(config.Jobs, jobs[v])
		}
	}*/
	return
}
