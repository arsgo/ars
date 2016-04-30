package mqservice

import (
	"encoding/json"
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

func NewStompService(sconfig string) IMQService {
	p := &StompService{}
	err := json.Unmarshal([]byte(sconfig), &p.config)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if strings.EqualFold(p.config.Address, "") {
		fmt.Println("address is nil")
		return nil
	}
	p.broker, err = stomp.NewStomp(p.config.Address)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return p
}
func (k *StompService) Send(queue string, msg string) (err error) {
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
