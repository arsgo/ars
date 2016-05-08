package cluster

import "encoding/json"

//WatchJobConfigChange 监控JOB配置变化
func (client *ClusterClient) WatchJobConfigChange(callback func(config *JobItems, err error)) {
	client.WaitClusterPathExists(client.jobConfigPath, client.timeout, func(exists bool) {
		if exists {
			callback(client.GetJobConfig(client.jobConfigPath))
		} else {
			client.Log.Info("job config path not exists")
		}
	})
	client.Log.Info("::watch job config changes")
	client.WatchClusterValueChange(client.jobConfigPath, func() {
		client.Log.Info("job config has changed")
		callback(client.GetJobConfig(client.jobConfigPath))
	})
}

//GetJobConfigs 获取JOB配置信息
func (client *ClusterClient) GetJobConfig(path string) (items *JobItems, err error) {
	items = &JobItems{}
	items.Jobs = make(map[string]JobItem)
	if !client.handler.Exists(path) {
		return
	}
	value, err := client.handler.GetValue(path)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(value), &items.Jobs)
	return
}

//GetJobConsumers 获取指定名称的JOBConsumer列表
func (client *ClusterClient) GetJobConsumers(jobName string) (jobs []string) {
	dmap := client.dataMap.Copy()
	dmap.Set("jobName", jobName)
	path := dmap.Translate(client.jobConsumerNamedRootFormat)
	children, err := client.handler.GetChildren(path)
	if err != nil {
		client.Log.Error(err)
		return
	}

	for _, v := range children {
		dmap.Set("path", v)
		values, err := client.handler.GetValue(dmap.Translate(client.jobConsumerRealPathFormat))
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
		jobs = append(jobs, consumer.Address)
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
