package cluster

import (
	"encoding/json"
	"strings"
)

//WatchRPCServiceChange 监控已发布的RPC服务变化
func (client *ClusterClient) WatchRPCServiceChange(callback func(services map[string][]string, err error)) {
	client.WaitClusterPathExists(client.rpcPublishPath, client.timeout, func(exists bool) {
		if !exists {
			client.Log.Info("service publish config not exists")
		} else {
			go func() {
				defer client.recover()
				callback(client.GetRPCService())
			}()
		}
	})
	client.Log.Info("::watch for service config changes ", client.rpcPublishPath)
	client.WatchClusterValueChange(client.rpcPublishPath, func() {
		go func() {
			defer client.recover()
			callback(client.GetRPCService())
		}()
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

//FilterRPCService 过滤RPC服务
func (client *ClusterClient) FilterRPCService(services map[string][]string) (items []TaskItem, err error) {
	all, err := client.GetServiceTasks()
	indentity := make(map[string]string)
	if err != nil {
		return
	}
	for _, v := range all.Tasks {
		v.Name = client.GetServiceFullPath(v.Name)
		if _, ok := indentity[v.Name]; !ok {
			indentity[v.Name] = v.Name
		}
		if _, ok := services[v.Name]; ok {
			items = append(items, v)
		}
	}
	for name := range services {
		if _, ok := indentity[name]; !ok {
			item := TaskItem{Name: name}
			items = append(items, item)
		}

	}
	return
}

//PublishRPCServices 发布所有服务
func (client *ClusterClient) PublishRPCServices(services ServiceProviderList) (err error) {
	client.publishLock.Lock()
	equal := services.Equal(client.lastServiceProviderList)
	client.lastServiceProviderList = services
	client.publishLock.Unlock()
	if equal {
		client.Log.Info("服务无变化")
		return
	}

	buffer, err := json.Marshal(services)
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

//GetServiceFullPath 获了服务的全名
func (client *ClusterClient) GetServiceFullPath(name string) string {
	if strings.Contains(name, "@") {
		return name
	}
	return name + client.domainPath
}
