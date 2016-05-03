package cluster

import (
	"errors"
	"fmt"

	"github.com/colinyl/lib4go/mq"
	"github.com/colinyl/lib4go/utility"
)

type msgHandler struct {
	service mq.IMQService
	Message string
}

func (h *msgHandler) Send(queue string, content string) error {
	return h.service.Send(queue, content)
}

type mqConsumer struct {
	service mq.IMQService
	param   map[string]interface{}
	handler func(*msgHandler) bool
}

func NewMQConsumer(param string, config string, handler func(*msgHandler) bool) (m *mqConsumer, err error) {
	m = &mqConsumer{}
	m.param, err = utility.GetParamsMap(param)
	if err != nil {
		return
	}
	if _, ok := m.param["queue"]; !ok {
		err = errors.New("queue is nil in params")
		return
	}
	m.service, err = mq.NewMQService(config)
	if err != nil {
		return
	}
	m.handler = handler
	if m.handler == nil {
		err = errors.New("handler is nill")
	}
	return
}
func (m *mqConsumer) Stop() {
	if m.service != nil {
		m.service.UnConsume(fmt.Sprintf("%s", m.param["queue"]))
		m.service.Close()
	}
}
func (m *mqConsumer) Start() {
	queue := fmt.Sprintf("%s", m.param["queue"])
	if m.service == nil {
		return
	}
	m.service.Consume(queue, func(h mq.MsgHandler) {
		msg := h.GetMessage()
		result := m.handler(&msgHandler{Message: msg, service: m.service})
		if result {
			h.Ack()
		} else {
			h.Nack()
		}
	})
}
