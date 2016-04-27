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
	Address string
}

func NewStompService(sconfig string) IMQService {
	p := &StompService{}
	err := json.Unmarshal([]byte(sconfig), &p.config)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if strings.EqualFold(p.config.Address, "") {		
		return nil
	}
	p.broker, err = stomp.NewStomp(p.config.Address)
	if err != nil {
		return nil
	}
	return p
}
func (k *StompService) Send(queue string, msg string) (err error) {
	return k.broker.Send(queue, msg)
}

func (k *StompService) Consume(queue string, callback func(stomp.MsgHandler)) (err error) {
	return k.broker.Consume(queue, callback)
}

func (k *StompService) Close() {
	k.broker.Close()
}
