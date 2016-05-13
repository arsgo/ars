package mqservice

import (
	"fmt"
	"strings"
	"sync"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/logger"
)

//MQHandler MQ任务处理程序
type MQHandler interface {
	Handle(cluster.TaskItem, string) bool
}

//MQConsumerService MQ消费服务
type MQConsumerService struct {
	clusterClient cluster.IClusterClient
	handler       MQHandler
	tasks         []cluster.TaskItem
	consumers     map[string]*MQConsumer
	lock          sync.Mutex
	Log           *logger.Logger
}

//NewMQConsumerService 创建MQ
func NewMQConsumerService(client cluster.IClusterClient, handler MQHandler) (mq *MQConsumerService, err error) {
	mq = &MQConsumerService{}
	mq.clusterClient = client
	mq.handler = handler
	mq.Log, err = logger.New("mq consumer", true)
	return
}

//UpdateTasks 更新MQ Consumer服务
func (mq *MQConsumerService) UpdateTasks(tasks []cluster.TaskItem) (err error) {
	consumers := mq.getTasks(tasks)
	mq.lock.Lock()
	defer mq.lock.Unlock()

	//关闭已启动的服务
	for k, v := range mq.consumers {
		if _, ok := consumers[k]; !ok {
			v.Stop()
			delete(mq.consumers, k)
		}
	}

	//启动已添加的服务
	for k, v := range consumers {
		if _, ok := mq.consumers[k]; !ok {
			mq.consumers[k], err = NewMQConsumer(v, mq.clusterClient, func(msg string) bool {
				return mq.handler.Handle(v, msg)
			})
			if err != nil {
				mq.Log.Fatal("mq create error:", err)
				continue
			}
			fmt.Println("::Start MQ Consumer:", v.Name)
			go mq.consumers[k].Start()
		}
	}
	return

}

//getTasks 获取当前服务列表
func (mq *MQConsumerService) getTasks(tasks []cluster.TaskItem) (consumers map[string]cluster.TaskItem) {
	consumers = make(map[string]cluster.TaskItem)
	for _, v := range tasks {
		if task, ok := consumers[v.Name]; !ok && strings.EqualFold(strings.ToLower(task.Type), "mq") && strings.EqualFold(strings.ToLower(task.Method), "consumer") {
			consumers[v.Name] = v
		}
	}
	return
}
