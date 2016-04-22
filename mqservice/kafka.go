package mqservice

import (
	"fmt"
	"strings"

	"github.com/jdamick/kafka"
)

type kafkaConfig struct {
	Address   string
	Topic     string
	Partition int
	Count     int
}

type KafkaPublisher struct {
	broker *kafka.BrokerPublisher
}
type KafkaConsumer struct {
	broker   *kafka.BrokerConsumer
	msgChan  chan *kafka.Message
	quitChan chan struct{}
}

func NewKafkaPublisher(config *kafkaConfig) (p *KafkaPublisher) {
	if strings.EqualFold(config.Address, "") ||
		strings.EqualFold(config.Topic, "") {
		fmt.Println("address or topic not allowed nil")
		return nil
	}
	p = &KafkaPublisher{}
	p.broker = kafka.NewBrokerPublisher(config.Address, config.Topic, 0)
	return
}
func (k *KafkaPublisher) Publish(content string) (err error) {
	_, err = k.broker.Publish(kafka.NewMessage([]byte(content)))
	return
}

func NewKafkaConsumer(config *kafkaConfig) (p *KafkaConsumer) {
	p.broker = kafka.NewBrokerConsumer(config.Address, config.Topic, 0, 0, 1048576)
	p.msgChan = make(chan *kafka.Message, config.Count)
	return
}
func (k *KafkaConsumer) Consume(callback func(string)) {
	go k.broker.ConsumeOnChannel(k.msgChan, 10, k.quitChan)
LOOP:
	for {
		select {
		case <-k.quitChan:
			break LOOP
		case msg := <-k.msgChan:
			callback(msg.PayloadString())
		}
	}
}
