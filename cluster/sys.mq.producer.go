package cluster

import "github.com/colinyl/lib4go/mq"

type MQProducer struct {
	service mq.IMQService
}

func (p *MQProducer) Send(queue string, content string) error {
	return p.service.Send(queue, content)
}

func NewMQProducer(config string) (m *MQProducer, err error) {
	m = &MQProducer{}
	m.service, err = mq.NewMQService(config)
	return
}
