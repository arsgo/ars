package mqservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/colinyl/stomp"
)

type StompService struct {
	config *StompConfig
	broker *stomp.Stomp
}

type StompConfig struct {
	Address string `json:"address"`
}

func NewStompService(sconfig string) (ps IMQService, err error) {
	p := &StompService{}
	ps = p
	err = json.Unmarshal([]byte(sconfig), &p.config)
	if err != nil {
		return
	}
	if strings.EqualFold(p.config.Address, "") {
		err = errors.New("address is nil")
		return
	}
	p.broker, err = stomp.NewStomp(p.config.Address)
	if err != nil {
		return
	}
	return
}
func (k *StompService) Send(queue string, msg string) (err error) {
	fmt.Printf("send:%s-%s", queue, msg)
	return k.broker.Send(queue, msg)
}

func (k *StompService) Consume(queue string, callback func(stomp.MsgHandler) bool) (err error) {
	return nil
	//return k.broker.Consume(queue, 10, callback)
}

func (k *StompService) Close() {
	k.broker.Close()
}
func StaticSend(queue string, msg string) (err error) {
	broker, err := stomp.NewStomp("192.168.101.161:61613")
	if err != nil {
		return
	}
	err = broker.Send(queue, msg)
	if err != nil {
		return
	}
	//	time.Sleep(time.Second)
	broker.Close()
	return
}
