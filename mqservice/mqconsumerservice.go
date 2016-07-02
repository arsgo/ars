package mqservice

import (
	"strings"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/concurrent"
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
	consumers     concurrent.ConcurrentMap //map[string]*MQConsumer
	Log           logger.ILogger
}

func (mq *MQConsumerService) recover() {
	if r := recover(); r != nil {
		mq.Log.Fatal(r)
	}
}

//NewMQConsumerService 创建MQ
func NewMQConsumerService(client cluster.IClusterClient, handler MQHandler, loggerName string) (mq *MQConsumerService, err error) {
	mq = &MQConsumerService{}
	mq.clusterClient = client
	mq.handler = handler
	mq.Log, err = logger.Get(loggerName, true)
	mq.consumers = concurrent.NewConcurrentMap()
	return
}

//UpdateTasks 更新MQ Consumer服务
func (mq *MQConsumerService) UpdateTasks(tasks []cluster.TaskItem) (err error) {
	consumers := mq.getTasks(tasks)
	//关闭已启动的服务
	currentConsumers := mq.consumers.GetAll()
	for k, v := range currentConsumers {
		if _, ok := consumers[k]; !ok {
			v.(*MQConsumer).Stop()
			mq.consumers.Delete(k)
		}
	}

	//启动已添加的服务
	for k, v := range consumers {	
		if c := mq.consumers.Get(k); c == nil {
			current, err := NewMQConsumer(v, mq.clusterClient, func(msg string, tk cluster.TaskItem) bool {
				return mq.handler.Handle(tk, msg)
			})
			if err != nil {
				mq.Log.Fatal("mq create error:", err)
				continue
			}
			mq.Log.Infof("::start mq consumer:[%s] %s", v.Name, v.Script)
			mq.consumers.Set(k, current)
			go func(mq *MQConsumerService, current *MQConsumer) {
				mq.recover()
				current.Start()
			}(mq, current)
		}
	}
	return

}

//getTasks 获取当前服务列表
func (mq *MQConsumerService) getTasks(tasks []cluster.TaskItem) (consumers map[string]cluster.TaskItem) {
	consumers = make(map[string]cluster.TaskItem)
	for _, v := range tasks {
		if _, ok := consumers[v.Name]; ok {
			continue
		}
		if strings.EqualFold(strings.ToLower(v.Type), "mq") && strings.EqualFold(strings.ToLower(v.Method), "consume") {
			consumers[v.Name] = v
		}
	}
	return
}
