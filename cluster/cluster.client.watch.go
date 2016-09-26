package cluster

import "time"

//WaitClusterPathExists  等待集群中的指定配置出现,不存在时持续等待
func (client *ClusterClient) WaitClusterPathExists(path string, timeout time.Duration, callback func(p string, exists bool)) {
	if p, ok := client.handler.Exists(path); ok {
		callback(p, true)
		return
	}
	callback("", false)
	timePiker := time.NewTicker(time.Second * 2)
	timeoutPiker := time.NewTicker(timeout)
	exists := false
	npath := ""
CHECKER:
	for {
		select {
		case <-timeoutPiker.C:
			break
		case <-timePiker.C:
			if v, ok := client.handler.Exists(path); ok {
				exists = true
				npath = v
				break CHECKER
			}
		}
	}
	callback(npath, exists)
}

//WatchClusterValueChange 等待集群指定路径的值的变化
func (client *ClusterClient) WatchClusterValueChange(path string, callback func()) {
	changes := make(chan string, 10)
	go func() {
		defer client.recover()
		client.handler.WatchValue(path, changes)
	}()
	go func() {
		for {
			select {
			case <-changes:
				callback()
			}
		}
	}()

}

//WatchClusterChildrenChange 监控集群指定路径的子节点变化
func (client *ClusterClient) WatchClusterChildrenChange(path string, callback func()) {
	changes := make(chan []string, 10)
	go func() {
		defer client.recover()
		client.handler.WatchChildren(path, changes)
	}()
	go func() {
		for {
			select {
			case <-changes:
				callback()
			}
		}
	}()
}
