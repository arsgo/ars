package mq

import (
	"fmt"
	"strings"

	"github.com/arsgo/ars/cluster"
	q "github.com/arsgo/lib4go/mq"
	"github.com/arsgo/lib4go/utility"
)

//MQConsumer MQ消费者
type MQConsumer struct {
	service q.IMQService
	param   map[string]interface{}
	handler func(string, cluster.TaskItem) bool
	queue   string
	setting string
	Message string
	task    cluster.TaskItem
}

//NewMQConsumer 创建新的MQ消费者
func NewMQConsumer(task cluster.TaskItem, clusterClient cluster.IClusterClient, handler func(string, cluster.TaskItem) bool) (mq *MQConsumer, err error) {
	mq = &MQConsumer{}
	mq.handler = handler
	mq.task = task
	mq.param, err = utility.GetParamsMap(task.Params)
	if err != nil {
		err = fmt.Errorf("mq consumer创建失败，获取mq参数失败:%s", task.Params)
		return
	}
	mq.queue = fmt.Sprintf("%s", mq.param["queue"])
	mq.setting = fmt.Sprintf("%s", mq.param["mq"])
	if strings.EqualFold(mq.queue, "") || strings.EqualFold(mq.setting, "") {
		err = fmt.Errorf("mq consumer创建失败，queue或setting不能为空：%s", task.Params)
		return
	}
	config, err := clusterClient.GetMQConfig(mq.setting)
	if err != nil {
		err = fmt.Errorf("mq consumer创建失败，配置文件格式有误：%s", mq.setting)
		return
	}
	mq.service, err = q.NewMQService(config)
	return
}

//Stop 停止服务
func (mq *MQConsumer) Stop() {
	if mq.service != nil {
		mq.service.UnConsume(mq.queue)
		mq.service.Close()
	}
}

//Start 启动服务
func (mq *MQConsumer) Start() {
	if mq.service == nil {
		return
	}
	mq.service.Consume(mq.queue, func(h q.MsgHandler) {
		msg := h.GetMessage()
		if mq.handler(msg, mq.task) {
			h.Ack()
		} else {
			h.Nack()
		}
	})
}
