package cluster

import "encoding/json"

//WatchAppTaskChange 监控APP Server的配置文件变化
func (client *ClusterClient) WatchAppTaskChange(callback func(config *AppServerTask, err error) error) {
	client.WaitClusterPathExists(client.appServerTaskPath, client.timeout, func(path string, exists bool) {
		if !exists {
			client.Log.Errorf("app config:%s未配置或不存在", client.appServerTaskPath)
		} else {
			go func() {
				defer client.recover()
				callback(client.GetCurrentAppServerTask())
			}()
		}
	})
	client.Log.Infof("::监控app config:%s的变化", client.appServerTaskPath)
	client.WatchClusterValueChange(client.appServerTaskPath, func() {
		client.Log.Infof(" -> app config:%s 值发生变化", client.appServerTaskPath)
		go func() {
			defer client.recover()
			callback(client.GetCurrentAppServerTask())
		}()
	})
}

//GetAppServerTaskNames 获取所有app server task名称
func (client *ClusterClient) GetAppServerTaskNames() (names []string, err error) {
	names, err = client.handler.GetChildren(client.appTaskRoot)
	return
}

//UpdateAppServerTask 更新AppServer配置文件
func (client *ClusterClient) UpdateAppServerTask(ip string, config *AppServerTask) (err error) {
	buffer, err := json.Marshal(config)
	if err != nil {
		client.Log.Errorf(" -> *AppServerTask转换为json出错：%v", config)
		return
	}
	nmap := client.dataMap.Copy()
	nmap.Set("ip", ip)
	path := nmap.Translate(p_appTaskConfig)
	err = client.handler.UpdateValue(path, string(buffer))
	return
}

//GetCurrentAppServerTask 获取当前AppServer任务
func (client *ClusterClient) GetCurrentAppServerTask() (config *AppServerTask, err error) {
	return client.GetAppServerTask(client.IP)
}

//GetAppServerTask 获取App Server的配置数据
func (client *ClusterClient) GetAppServerTask(ip string) (config *AppServerTask, err error) {
	nmap := client.dataMap.Copy()
	nmap.Set("ip", ip)
	path := nmap.Translate(p_appTaskConfig)
	config = &AppServerTask{}
	values, err := client.handler.GetValue(path)
	if err != nil {
		client.Log.Errorf(" -> app config：%s 获取配置数据有误", path)
		return
	}
	err = json.Unmarshal([]byte(values), &config)
	if err != nil {
		client.Log.Errorf(" -> app config：%s json格式有误", path)
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
func (client *ClusterClient) CreateAppServer(port string, snap string) (path string, err error) {
	data := client.dataMap.Copy()
	data.Set("ip", client.IP)
	data.Set("port", port)
	path = data.Translate(p_appServerPath)
	err = client.SetNode(path, snap)
	return
}
func (d *ClusterClient) CloseAppServer(path string) (err error) {
	return d.CloseNode(path)
}

//CreateMQConsumer 创建mq conusmer
func (client *ClusterClient) CreateMQConsumer(name string, value string) (string, error) {
	data := client.dataMap.Copy()
	data.Set("name", name)
	path := data.Translate(p_MQConsumerClusterClientBase)
	return client.handler.CreateSeqNode(path, value)
}
