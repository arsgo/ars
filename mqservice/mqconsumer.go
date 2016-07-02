package mqservice

import (
	"errors"
	"fmt"
	"strings"

	"github.com/colinyl/ars/cluster"
	queue "github.com/colinyl/lib4go/mq"
	"github.com/colinyl/lib4go/utility"
)

//MQConsumer MQ消费者
type MQConsumer struct {
	service queue.IMQService
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
		return
	}
	mq.queue = fmt.Sprintf("%s", mq.param["queue"])
	mq.setting = fmt.Sprintf("%s", mq.param["mq"])
	if strings.EqualFold(mq.queue, "") || strings.EqualFold(mq.setting, "") {
		err = errors.New("queue name  or mq name  is nil in params")
		return
	}
	config, err := clusterClient.GetMQConfig(mq.setting)
	if err != nil {
		return
	}
	mq.service, err = queue.NewMQService(config)
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
	mq.service.Consume(mq.queue, func(h queue.MsgHandler) {
		msg := h.GetMessage()
		if mq.handler(msg, mq.task) {
			h.Ack()
		} else {
			h.Nack()
		}
	})
}
