package main

import "github.com/arsgo/ars/cluster"

//OnMQConsumerCreate  创建JOB服务
func (a *SPServer) OnMQConsumerCreate(task cluster.TaskItem) (path string) {
	return ""
}

//OnMQConsumerClose 关闭JOB服务
func (a *SPServer) OnMQConsumerClose(task cluster.TaskItem, path string) {
}
