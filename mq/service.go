package mq

import (
	"runtime/debug"
	"strings"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/lib4go/concurrent"
	"github.com/arsgo/lib4go/logger"
	"github.com/arsgo/lib4go/utility"
)

//MQHandler MQ任务处理程序
type MQHandler interface {
	Handle(cluster.TaskItem, string, string) bool
	OnOpenTask(task cluster.TaskItem) string
	OnCloseTask(task cluster.TaskItem, path string)
}

//MQConsumerService MQ消费服务
type MQConsumerService struct {
	clusterClient cluster.IClusterClient
	handler       MQHandler
	collector     base.ICollector
	tasks         []cluster.TaskItem
	Available     bool
	paths         *concurrent.ConcurrentMap
	consumers     *concurrent.ConcurrentMap //map[string]*MQConsumer
	Log           logger.ILogger
}

func (mq *MQConsumerService) recover() {
	if r := recover(); r != nil {
		mq.Log.Fatal(r, string(debug.Stack()))
	}
}

//NewMQConsumerService 创建MQ
func NewMQConsumerService(client cluster.IClusterClient, handler MQHandler, loggerName string, collector base.ICollector) (mq *MQConsumerService, err error) {
	mq = &MQConsumerService{clusterClient: client, handler: handler, collector: collector, Available: false}
	mq.Log, err = logger.Get(loggerName)
	mq.consumers = concurrent.NewConcurrentMap()
	mq.paths = concurrent.NewConcurrentMap()
	return
}
func (mq *MQConsumerService) GetServices() (svs map[string]string) {
	svs = make(map[string]string)
	services := mq.paths.GetAll()
	for i, v := range services {
		svs[i] = v.(string)
	}
	return

}

//UpdateTasks 更新MQ Consumer服务
func (mq *MQConsumerService) UpdateTasks(tasks []cluster.TaskItem) (err error) {
	consumers := mq.getTasks(tasks)
	//关闭已启动的服务
	currentConsumers := mq.consumers.GetAll()
	for k, v := range currentConsumers {
		if tks, ok := consumers[k]; !ok {
			mq.Log.Info(" -> 关闭 mq consumer:", k)
			v.(*MQConsumer).Stop()
			mq.consumers.Delete(k)
			if mq.paths.Get(k) != nil {
				mq.handler.OnCloseTask(tks, mq.paths.Get(k).(string))
				mq.paths.Delete(k)
			}
		}
	}

	//启动已添加的服务
	for k, v := range consumers {
		if ok, _, _ := mq.consumers.Add(k, mq.createConsumer, v); ok {
			path := mq.handler.OnOpenTask(v)
			mq.paths.Set(v.Name, path)
		}
	}
	mq.Available = mq.consumers.GetLength() > 0
	return
}
func (mq *MQConsumerService) createConsumer(args ...interface{}) (r interface{}, err error) {
	v := args[0].(cluster.TaskItem)
	current, err := NewMQConsumer(v, mq.clusterClient, func(msg string, tk cluster.TaskItem) bool {
		return mq.handler.Handle(tk, msg, utility.GetSessionID())
	}, mq.collector)
	if err != nil {
		mq.Log.Error(" -> mq consumer 启动失败:", err)
		return
	}
	go func(mq *MQConsumerService, current *MQConsumer) {
		defer mq.recover()
		current.Start()
	}(mq, current)
	r = current
	return
}

//getTasks 获取当前服务列表
func (mq *MQConsumerService) getTasks(tasks []cluster.TaskItem) (consumers map[string]cluster.TaskItem) {
	consumers = make(map[string]cluster.TaskItem)
	for _, v := range tasks {
		if _, ok := consumers[v.Name]; ok {
			continue
		}
		if strings.EqualFold(v.Type, "mq") && strings.EqualFold(v.Method, "consume") {
			consumers[v.Name] = v
		}
	}
	return
}
