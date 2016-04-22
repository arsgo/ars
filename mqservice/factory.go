package mqservice

import (
	"encoding/json"
	"fmt"
)

const (
	stompMQ = "stomp"
)

type MQConfig struct {
	Type string
}

func NewMQService(config string) IMQService {
	p := &MQConfig{}
	err := json.Unmarshal([]byte(config), &p)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	switch p.Type {
	case stompMQ:	
		return NewStompService(config)
	}

	return nil
}
