package cluster

import (
	"encoding/json"
	"fmt"
)

//WatchRCTaskChange 监控RC Config变化
func (client *ClusterClient) WatchRCTaskChange(callback func(RCServerTask, error)) {
	client.WaitClusterPathExists(client.rcServerConfig, client.timeout, func(path string, exists bool) {
		if !exists {
			client.Log.Errorf("rc config:%s未配置或不存在", client.rcServerConfig)
		} else {
			go func() {
				defer client.recover()
				callback(client.GetRCServerTask())
			}()
		}
	})
	client.Log.Infof("::监控rc config:%s的变化", client.rcServerConfig)
	client.WatchClusterValueChange(client.rcServerConfig, func() {
		client.Log.Infof(" -> rc config:%s 值发生变化", client.rcServerConfig)
		go func() {
			defer client.recover()
			callback(client.GetRCServerTask())
		}()
	})
}

//WatchRCServerChange 监控RC服务器变化,变化后回调指定函数
func (client *ClusterClient) WatchRCServerChange(callback func([]*RCServerItem, error)) {
	client.WaitClusterPathExists(client.rcServerRoot, client.timeout, func(path string, exists bool) {
		if !exists {
			client.Log.Errorf("rc servers:%s未配置或不存在", client.rcServerRoot)
		} else {
			go func() {
				defer client.recover()
				callback(client.GetAllRCServers())
			}()
		}
	})
	client.Log.Infof("::监控rc servers:%s的变化", client.rcServerRoot)
	client.WatchClusterChildrenChange(client.rcServerRoot, func() {
		client.Log.Infof(" -> rc servers:%s 值发生变化", client.rcServerRoot)
		go func() {
			defer client.recover()
			callback(client.GetAllRCServers())
		}()
	})
}

//UpdateRCServerTask 更新RcServer任务配置
func (client *ClusterClient) UpdateRCServerTask(config RCServerTask) (err error) {
	buffer, err := json.Marshal(config)
	if err != nil {
		client.Log.Errorf(" -> RCServerTask转换为json出错：%v", config)
		return
	}
	err = client.handler.UpdateValue(client.rcServerConfig, string(buffer))
	return
}

//GetRCServerValue 获取RC服务器信息
func (client *ClusterClient) GetRCServerValue(path string) (value *RCServerItem, err error) {
	content, err := client.handler.GetValue(path)
	if err != nil {
		client.Log.Errorf(" -> rc server:%s 获取server数据有误", path)
		return
	}
	value = &RCServerItem{}
	err = json.Unmarshal([]byte(content), &value)
	if err != nil {
		client.Log.Errorf(" -> rc server:%s json格式有误", content)
		return
	}
	value.Path = path
	return
}

//GetAllRCServers 获取所有RC服务器信息
func (client *ClusterClient) GetAllRCServers() (servers []*RCServerItem, err error) {
	rcs, err := client.handler.GetChildren(client.rcServerRoot)
	if err != nil {
		client.Log.Errorf(" -> 获取所有rc servers 出错:%s,%v", client.rcServerRoot, err)
		return
	}
	servers = []*RCServerItem{}
	for _, v := range rcs {
		rcPath := fmt.Sprintf("%s/%s", client.rcServerRoot, v)
		config, err := client.GetRCServerValue(rcPath)
		if err != nil {
			client.Log.Errorf(" -> 获取rc server数据有误:%v", err)
			continue
		}
		if len(config.Address) > 0 {
			servers = append(servers, config)
		}
	}
	return
}

//CreateRCServer 创建RCServer
func (client *ClusterClient) CreateRCServer(value string) (string, error) {
	return client.handler.CreateSeqNode(client.dataMap.Translate(p_rcServerClusterClientBase), value)
}

//CloseRCServer close rc server
func (client *ClusterClient) CloseRCServer(path string) error {
	return client.CloseNode(path)
}

//GetRCServerTask 获取RC Server任务
func (client *ClusterClient) GetRCServerTask() (config RCServerTask, err error) {
	value, err := client.handler.GetValue(client.rcServerConfig)
	if err != nil {
		client.Log.Errorf(" -> rc config：%s 获取配置数据有误", client.rcServerConfig)
		return
	}
	config = RCServerTask{}
	err = json.Unmarshal([]byte(value), &config)
	if err != nil {
		client.Log.Errorf(" -> rc config：%s 配置数据json格式有误，%v", client.rcServerConfig, err)
	}
	return
}
