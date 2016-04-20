package mqservice

import "encoding/json"

const (
	KafkaMQName = "kafka"
)

type MQConfig struct {
	Type string
}

func NewMQService(config string) IMQService {
	p := &MQConfig{}
	err := json.Unmarshal([]byte(config), &p)
	if err != nil {
		return nil
	}
	switch p.Type {
	case KafkaMQName:
		return NewKafkaService(config)
	}

	return nil
}
