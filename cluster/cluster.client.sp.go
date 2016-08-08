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
			client.Log.Errorf("sp config:%s未配置或不存在", client.spConfigPath)
		} else {
			go func() {
				defer client.recover()
				callback()
			}()
		}
	})
	client.Log.Infof("::监控sp config:%s的变化", client.spConfigPath)
	client.WatchClusterValueChange(client.spConfigPath, func() {
		client.Log.Infof(" -> sp config:%s 值发生变化", client.spConfigPath)
		go func() {
			defer client.recover()
			callback()
		}()
	})
}

//WatchSPServerChange 监控sp server变化
func (client *ClusterClient) WatchSPServerChange(changed func(RPCServices, error)) (err error) {

	client.WaitClusterPathExists(client.rpcProviderRootPath, client.timeout, func(exists bool) {
		if !exists {
			client.Log.Errorf("sp servers:%s未配置或不存在", client.rpcProviderRootPath)
		} else {
			go func() {
				defer client.recover()
				changed(client.GetSPServerServices())
			}()
		}
	})
	client.Log.Infof("::监控sp servers:%s的变化", client.rpcProviderRootPath)
	client.WatchClusterChildrenChange(client.rpcProviderRootPath, func() {
		go func() {
			defer client.recover()
			client.Log.Infof(" -> sp servers:%s 值发生变化", client.rpcProviderRootPath)
			changed(client.GetSPServerServices())
		}()
	})
	lst, err := client.GetAllSPServers()
	for _, v := range lst {
		for _, p := range v {
			client.WatchClusterChildrenChange(p, func() {
				go func() {
					defer client.recover()
					client.Log.Infof(" -> sp servers:%s 值发生变化", client.rpcProviderRootPath)
					changed(client.GetSPServerServices())
				}()
			})
		}
	}
	return
}

//GetAllSPServers 获了所有服务名称路径
func (client *ClusterClient) GetAllSPServers() (lst map[string][]string, err error) {
	lst = make(map[string][]string)
	serviceList, err := client.handler.GetChildren(client.rpcProviderRootPath)
	if err != nil {
		client.Log.Errorf(" -> sp server:%s 获取all servers数据有误", client.rpcProviderRootPath)
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

//GetSPServerServices 根据服务提供方路径,获取所有服务列表
func (client *ClusterClient) GetSPServerServices() (lst RPCServices, err error) {
	lst = make(map[string][]string)
	serviceList, err := client.handler.GetChildren(client.rpcProviderRootPath)
	if err != nil {
		client.Log.Errorf(" -> sp server:%s 获取children数据有误", client.rpcProviderRootPath)
		return
	}

	for _, value := range serviceList {
		name := client.GetServiceFullPath(value)
		path := fmt.Sprintf("%s/%s", client.rpcProviderRootPath, value)
		providerList, er := client.handler.GetChildren(path)
		if er != nil {
			client.Log.Errorf(" -> sp server:%s 获取children数据有误", path)
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
func (client *ClusterClient) UpdateSPServerTask(config SPServerTask) (err error) {
	buffer, err := json.Marshal(config)
	if err != nil {
		client.Log.Errorf(" -> SPServerTask转换为json出错：%v", config)
		return
	}
	err = client.handler.UpdateValue(client.spServerTaskPath, string(buffer))
	return
}

//GetSPServerTask 获取service provider 的任务列表
func (client *ClusterClient) GetSPServerTask(ip string) (task SPServerTask, err error) {

	task = SPServerTask{}
	values, err := client.handler.GetValue(client.spServerTaskPath)
	if err != nil {
		client.Log.Errorf(" -> sp config:%s 获取数据有误", client.spServerTaskPath)
		return
	}
	err = json.Unmarshal([]byte(values), &task)
	if err != nil {
		client.Log.Errorf(" -> sp config：%s json格式有误", values)
		return
	}
	var items []TaskItem
	for _, v := range task.Tasks {
		if strings.EqualFold(ip, "*") || strings.EqualFold(v.IP, "*") || strings.Contains(","+v.IP+",", ip) {
			v.Name = client.GetServiceFullPath(v.Name)
			items = append(items, v)
		}
	}
	task.Tasks = items
	return
}

//CreateSPServer 创建服务提供节点
func (client *ClusterClient) CreateSPServer(name string, port string, value string) (string, error) {
	data := client.dataMap.Copy()
	data.Set("serviceName", strings.TrimSuffix(name, client.domainPath))
	data.Set("ip", client.IP)
	data.Set("port", port)
	path := data.Translate(p_serviceProviderPath)
	return client.handler.CreateTmpNode(path, value)

}

//CloseSPServer 关闭sp server节点
func (client *ClusterClient) CloseSPServer(path string) error {
	return client.CloseNode(path)
}
