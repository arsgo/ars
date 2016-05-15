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
			go callback(client.GetRPCService())
		}
	})
	client.Log.Info("::watch for service config changes ")
	client.WatchClusterValueChange(client.rpcPublishPath, func() {
		client.Log.Info(" -> rpc serivce has changed")
		go callback(client.GetRPCService())
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
	if err != nil {
		return
	}
	for _, v := range all.Tasks {
		v.Name = client.GetServiceFullPath(v.Name)
		if _, ok := services[v.Name]; ok {
			items = append(items, v)
		}
	}
	return
}

//PublishRPCServices 发布所有服务
func (client *ClusterClient) PublishRPCServices(crossServices map[string]map[string][]string) (err error) {
	crossServices = make(map[string]map[string][]string)
	providers, err := client.GetServiceProviderPaths()
	if err != nil {
		return
	}
	//处理跨域服务
	if crossServices != nil {
		for domain, services := range crossServices {
			for service, ips := range services {
				providers[service+"@"+domain] = ips
			}
		}
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

//GetServiceFullPath 获了服务的全名
func (client *ClusterClient) GetServiceFullPath(name string) string {
	if strings.Contains(name, "@") {
		return name
	}
	return name + client.domainPath
}
