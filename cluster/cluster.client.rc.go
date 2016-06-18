package cluster

import (
	"encoding/json"
	"fmt"
)

//WatchRCServerChange 监控RC服务器变化,变化后回调指定函数
func (client *ClusterClient) WatchRCServerChange(callback func([]*RCServerItem, error)) {
	client.WaitClusterPathExists(client.rcServerRoot, client.timeout, func(exists bool) {
		if !exists {
			client.Log.Infof("rc server not exists:%s", client.rcServerRoot)
		} else {
			go callback(client.GetAllRCServerValues())
		}
	})
	client.Log.Info("::watch for rc server changes")
	client.WatchClusterChildrenChange(client.rcServerRoot, func() {
		client.Log.Info(" -> rc server has changed")
		go callback(client.GetAllRCServerValues())
	})
}

//GetRCServerValue 获取RC服务器信息
func (client *ClusterClient) GetRCServerValue(path string) (value *RCServerItem, err error) {
	content, err := client.handler.GetValue(path)
	if err != nil {
		return
	}
	value = &RCServerItem{}
	err = json.Unmarshal([]byte(content), &value)
	value.Path = path
	return
}

//GetAllRCServerValues 获取所有RC服务器信息
func (client *ClusterClient) GetAllRCServerValues() (servers []*RCServerItem, err error) {
	rcs, _ := client.handler.GetChildren(client.rcServerRoot)
	servers = []*RCServerItem{}
	for _, v := range rcs {
		rcPath := fmt.Sprintf("%s/%s", client.rcServerRoot, v)
		config, err := client.GetRCServerValue(rcPath)
		if err != nil {
			continue
		}
		servers = append(servers, config)
	}
	return
}

//CreateRCServer 创建RCServer
func (client *ClusterClient) CreateRCServer(value string) (string, error) {
	return client.handler.CreateSeqNode(client.dataMap.Translate(p_rcServerClusterClientBase), value)
}

//GetRCServerTasks 获取RC Server任务
func (client *ClusterClient) GetRCServerTasks() (config RCServerTask, err error) {
	value, err := client.handler.GetValue(client.rcServerConfig)
	if err != nil {
		return
	}
	config = RCServerTask{}
	err = json.Unmarshal([]byte(value), &config)
	return
}

//WatchRCTaskChange 监控RC Config变化
func (client *ClusterClient) WatchRCTaskChange(callback func(RCServerTask, error)) {
	client.WaitClusterPathExists(client.rcServerConfig, client.timeout, func(exists bool) {
		if !exists {
			client.Log.Infof("rc server config not exists:%s", client.rcServerRoot)
		} else {
			go callback(client.GetRCServerTasks())
		}
	})
	client.Log.Info("::watch for rc server config changes")
	client.WatchClusterValueChange(client.rcServerConfig, func() {
		client.Log.Info(" -> rc server config has changed")
		go callback(client.GetRCServerTasks())
	})
}
