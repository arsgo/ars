package cluster

import "encoding/json"

//WatchJobConfigChange 监控JOB配置变化
func (client *ClusterClient) WatchJobConfigChange(callback func(config map[string]JobItem, err error)) {
	client.WaitClusterPathExists(client.jobConfigPath, client.timeout, func(exists bool) {
		if exists {
			go func() {
				defer client.recover()
				callback(client.GetJobTask())
			}()
		} else {
			client.Log.Info("job config path not exists:", client.jobConfigPath)
		}
	})
	client.Log.Info("::watch for job config changes")
	client.WatchClusterValueChange(client.jobConfigPath, func() {
		client.Log.Info(" -> job config has changed")
		callback(client.GetJobTask())
	})
}

//GetJobTask 获取JOB配置信息
func (client *ClusterClient) GetJobTask() (items map[string]JobItem, err error) {
	path := client.jobConfigPath
	if !client.handler.Exists(path) {
		return
	}
	value, err := client.handler.GetValue(path)
	if err != nil {
		return
	}
	jobs := []JobItem{}
	items = make(map[string]JobItem)
	err = json.Unmarshal([]byte(value), &jobs)
	if err != nil {
		return
	}
	for _, v := range jobs {
		items[v.Name] = v
	}
	return
}

//UpdateJobTask 更新JOB配置信息
func (client *ClusterClient) UpdateJobTask(jobName string, config map[string]JobItem) (err error) {
	buffer, err := json.Marshal(config)
	if err != nil {
		return
	}
	err = client.handler.UpdateValue(client.jobConfigPath, string(buffer))
	return
}

//GetJobConsumers 获取指定名称的JOBConsumer列表
func (client *ClusterClient) GetJobConsumers(jobName string) (jobs []string) {
	dmap := client.dataMap.Copy()
	dmap.Set("jobName", jobName)
	root := dmap.Translate(p_jobConsumerNamedRootForamt)
	children, err := client.handler.GetChildren(root)
	if err != nil {
		client.Log.Error(err)
		return
	}
	for _, v := range children {
		dmap.Set("path", v)
		path := dmap.Translate(p_jobConsumerClusterClientPathFormat)
		values, err := client.handler.GetValue(path)
		if err != nil {
			client.Log.Error(err)
			continue
		}
		consumer := &JobConsumerValue{}
		err = json.Unmarshal([]byte(values), &consumer)
		if err != nil {
			client.Log.Error(err)
			continue
		}
		jobs = append(jobs, consumer.Server)
	}
	return
}
func (client *ClusterClient) UpdateJobConsumerPath(path string, value string) error {
	return client.handler.UpdateValue(path, value)

}
func (client *ClusterClient) CreateJobConsumer(jobName string, value string) (string, error) {
	data := client.dataMap.Copy()
	data.Set("jobName", jobName)
	path := data.Translate(p_jobConsumerClusterClientBase)
	return client.handler.CreateSeqNode(path, value)

}

func (client *ClusterClient) CloseJobConsumer(path string) error {
	return client.handler.Delete(path)
}
