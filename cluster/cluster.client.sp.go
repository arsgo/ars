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
			go func() {
				defer client.recover()
				callback()
			}()
		}
	})
	client.Log.Info("::watch for provider config change")
	client.WatchClusterValueChange(client.spConfigPath, func() {
		client.Log.Info(" -> service provider task has changed")
		go func() {
			defer client.recover()
			callback()
		}()
	})
}

//WatchServiceProviderChange 监控RPC服务提供方变化
func (client *ClusterClient) WatchServiceProviderChange(changed func(ServiceProviderList, error)) (err error) {
	client.Log.Info("::watch for service providers changes")
	client.WaitClusterPathExists(client.rpcProviderRootPath, client.timeout, func(exists bool) {
		if !exists {
			client.Log.Info("service provider node not exists:", client.rpcProviderRootPath)
		} else {
			go func() {
				defer client.recover()
				changed(client.GetServiceProviders())
			}()
		}
	})
	client.WatchClusterChildrenChange(client.rpcProviderRootPath, func() {
		go func() {
			defer client.recover()
			changed(client.GetServiceProviders())
		}()
	})
	lst, err := client.GetAllServiceProviderNamePath()
	for _, v := range lst {
		for _, p := range v {
			client.WatchClusterChildrenChange(p, func() {
				go func() {
					defer client.recover()
					changed(client.GetServiceProviders())
				}()
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

//GetServiceProviders 根据服务提供方路径,获取所有服务列表
func (client *ClusterClient) GetServiceProviders() (lst ServiceProviderList, err error) {
	lst = make(map[string][]string)
	serviceList, err := client.handler.GetChildren(client.rpcProviderRootPath)
	if err != nil {
		return
	}

	for _, value := range serviceList {
		name := client.GetServiceFullPath(value)
		path := fmt.Sprintf("%s/%s", client.rpcProviderRootPath, value)
		providerList, er := client.handler.GetChildren(path)
		if er != nil {
			continue
		}
		for _, l := range providerList {
			if _, ok := lst[name]; !ok {
				lst[name] = []string{}
			}
			lst[name] = append(lst[name], l)
		}
	}
	return
}

//UpdateSPServerTask 更新sp server task config
func (client *ClusterClient) UpdateSPServerTask(config ServiceProviderTask) (err error) {
	buffer, err := json.Marshal(config)
	if err != nil {
		return
	}
	err = client.handler.UpdateValue(client.spServerTaskPath, string(buffer))
	return
}

//GetSPServerTask 获取service provider 的任务列表
func (client *ClusterClient) GetSPServerTask() (task ServiceProviderTask, err error) {

	task = ServiceProviderTask{}
	values, err := client.handler.GetValue(client.spServerTaskPath)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(values), &task)
	if err != nil {
		return
	}
	var items []TaskItem
	for _, v := range task.Tasks {
		if strings.EqualFold(v.IP, "*") || strings.Contains(","+v.IP+",", client.IP) {
			v.Name = client.GetServiceFullPath(v.Name)
			items = append(items, v)
		}
	}
	task.Tasks = items
	return
}

//CreateServiceProvider 创建服务提供节点
func (client *ClusterClient) CreateServiceProvider(name string, port string, value string) (string, error) {
	data := client.dataMap.Copy()
	data.Set("serviceName", strings.TrimSuffix(name, client.domainPath))
	data.Set("ip", client.IP)
	data.Set("port", port)
	path := data.Translate(p_serviceProviderPath)
	return client.handler.CreateTmpNode(path, value)

}

func (client *ClusterClient) CloseServiceProvider(path string) error {
	return client.handler.Delete(path)
}
