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
			callback(client.GetAllRCServerValues())
		}
	})
	client.Log.Info("::watch for rc server changes")
	client.WatchClusterChildrenChange(client.rcServerRoot, func() {
		client.Log.Info("rc server has changed")
		callback(client.GetAllRCServerValues())
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
	return
}

//GetAllRCServers 获取所有RC服务器信息
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
func (client *ClusterClient) CreateRCServer(value string) (string, error) {
	return client.handler.CreateSeqNode(client.dataMap.Translate(p_rcServerClusterClientBase), value)
}
