package cluster

import (
	"encoding/json"
	"strings"
)

//WatchRPCServiceChange 监控已发布的RPC服务变化
func (client *ClusterClient) WatchRPCServiceChange(callback func(services map[string][]string, err error)) {
	client.WaitClusterPathExists(client.rpcPublishPath, client.timeout, func(exists bool) {
		if !exists {
			client.Log.Errorf("services:%s未配置或不存在", client.rpcPublishPath)
		} else {
			go func() {
				defer client.recover()
				callback(client.GetPublishServices())
			}()
		}
	})
	client.Log.Infof("::监控services:%s的变化", client.rpcPublishPath)
	client.WatchClusterValueChange(client.rpcPublishPath, func() {
		go func() {
			defer client.recover()
			client.Log.Infof(" -> services:%s 值发生变化", client.rpcPublishPath)
			callback(client.GetPublishServices())
		}()
	})
}

//GetPublishServices 获取已发布的RPC服务
func (client *ClusterClient) GetPublishServices() (sp RPCServices, err error) {
	content, err := client.handler.GetValue(client.rpcPublishPath)
	if err != nil {
		client.Log.Errorf(" -> services:%s 获取server数据有误", client.rpcPublishPath)
		return
	}
	sp = make(map[string][]string)
	err = json.Unmarshal([]byte(content), &sp)
	return
}

//GetLocalServices 过滤RPC服务
func (client *ClusterClient) GetLocalServices(services map[string][]string) (items []TaskItem, err error) {
	all, err := client.GetSPServerTask("*")
	if err != nil {
		return
	}
	indentity := make(map[string]string)
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

//PublishServices 发布所有服务
func (client *ClusterClient) PublishServices(services RPCServices) (err error) {
	client.publishLock.Lock()
	equal := services.Equal(client.lastRPCServices)
	client.lastRPCServices = services.Clone()
	client.publishLock.Unlock()
	if equal {
		client.Log.Debug("服务无变化")
		return
	}

	buffer, err := json.Marshal(services)
	if err != nil {
		client.Log.Errorf(" -> services转换为json出错：%v", services)
		return
	}
	err = client.SetNode(client.rpcPublishPath, string(buffer))
	return
}

//GetServiceFullPath 获了服务的全名
func (client *ClusterClient) GetServiceFullPath(name string) string {
	if strings.Contains(name, "@") {
		return name
	}
	return name + client.domainPath
}
