package cluster

import "encoding/json"

//WatchRPCServiceChange 监控已发布的RPC服务变化
func (client *ClusterClient) WatchRPCServiceChange(callback func(services map[string][]string, err error)) {
	client.WaitClusterPathExists(client.rpcPublishPath, client.timeout, func(exists bool) {
		if !exists {
			client.Log.Info("service publish config not exists")
		} else {
			callback(client.GetRPCService())
		}
	})
	client.Log.Info("::watch for service config changes ")
	client.WatchClusterValueChange(client.rpcPublishPath, func() {
		client.Log.Info("serivce has changed")
		callback(client.GetRPCService())
	})
}

//GetRPCService 获取已发布的RPC服务
func (client *ClusterClient) GetRPCService() (sp ServiceProviderList, err error) {
	content, err := client.handler.GetValue(client.rpcPublishPath)
	if err != nil {
		return
	}
	sp = make(map[string][]string)
	err = json.Unmarshal([]byte(content), &sp)
	return
}


//PublishRPCServices 发布所有服务
func (client *ClusterClient) PublishRPCServices() (err error) {

	providers, err := client.GetServiceProviderPaths()
	if err != nil {
		return
	}
	buffer, err := json.Marshal(providers)
	if err != nil {
		return
	}
	serviceValue := string(buffer)
	if client.handler.Exists(client.rpcPublishPath) {
		err = client.handler.UpdateValue(client.rpcPublishPath, serviceValue)
	} else {
		err = client.handler.CreatePath(client.rpcPublishPath, serviceValue)
	}
	return
}