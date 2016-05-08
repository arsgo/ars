package cluster

import (
	"fmt"
	"strings"
)
import "encoding/json"

//WatchSPTaskChange 监控SP Task任务变化
func (client *ClusterClient) WatchSPTaskChange(callback func()) {
	client.WaitClusterPathExists(client.spConfigPath, client.timeout, func(exists bool) {
		if !exists {
			client.Log.Info("sp server config not exists")

		} else {
			callback()
		}
	})
	client.Log.Info("::watch for provider config change")
	client.WatchClusterValueChange(client.spConfigPath, func() {
		callback()
	})
}

//WatchServiceProviderChange 监控RPC服务提供方变化
func (client *ClusterClient) WatchServiceProviderChange(changed func()) (err error) {

	client.Log.Info("::watch for service providers changes")
	client.WaitClusterPathExists(client.rpcProviderRootPath, client.timeout, func(exists bool) {
		if !exists {
			client.Log.Info("service provider node not exists")
		} else {
			err = client.PublishRPCServices(nil)
			changed()
		}
	})
	client.WatchClusterChildrenChange(client.rpcProviderRootPath, func() {
		err = client.PublishRPCServices(nil)
		changed()
	})
	lst, err := client.GetAllServiceProviderNamePath()
	for _, v := range lst {
		for _, p := range v {
			client.WatchClusterChildrenChange(p, func() {
				err = client.PublishRPCServices(nil)
				changed()
			})
		}

	}
	return
}

//GetAllServiceProviderNamePath 获了所有服务名称路径
func (client *ClusterClient) GetAllServiceProviderNamePath() (lst map[string][]string, err error) {
	lst = make(map[string][]string)
	serviceList, err := client.handler.GetChildren(client.rpcProviderRootPath)
	if err != nil {
		return
	}
	for _, v := range serviceList {
		if _, ok := lst[v]; !ok {
			lst[v] = []string{}
		}
		lst[v] = append(lst[v], fmt.Sprintf("%s/%s", client.rpcProviderRootPath, v))
	}
	return
}

//GetServiceProviderPaths 根据服务提供方路径,获取所有服务列表
func (client *ClusterClient) GetServiceProviderPaths() (lst ServiceProviderList, err error) {
	lst = make(map[string][]string)
	serviceList, err := client.handler.GetChildren(client.rpcProviderRootPath)
	if err != nil {
		return
	}

	for _, v := range serviceList {
		path := fmt.Sprintf("%s/%s", client.rpcProviderRootPath, v)
		providerList, er := client.handler.GetChildren(path)
		if er != nil {
			continue
		}
		for _, l := range providerList {
			if _, ok := lst[v]; !ok {
				lst[v] = []string{}
			}
			lst[v] = append(lst[v], l)
		}
	}
	return
}

//GeServiceTasks 获取service provider 的任务列表
func (client *ClusterClient) GeServiceTasks() (items []TaskItem, err error) {
	var taskItem []TaskItem
	values, err := client.handler.GetValue(client.spServerTaskPath)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(values), &taskItem)
	if err != nil {
		return
	}
	for _, v := range taskItem {
		if strings.EqualFold(v.IP, "*") ||
			strings.Contains(","+v.IP+",", client.IP) {
			items = append(items, v)
		}
	}

	return
}

//CreateServiceProvider 创建服务提供节点
func (client *ClusterClient) CreateServiceProvider(name string, port string, value string) (string, error) {
	data := client.dataMap.Copy()
	data.Set("serviceName", name)
	data.Set("ip", client.IP)
	data.Set("port", port)
	path := data.Translate(p_serviceProviderPath)
	return client.handler.CreateTmpNode(path, value)

}
func (client *ClusterClient) CloseServiceProvider(path string) error {
	return client.handler.Delete(path)
}