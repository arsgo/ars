package cluster

import (
	"fmt"
	"strings"
	"sync"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

type MQConsumerHandler interface {
	Hande(string, string, *msgHandler) bool
	GetSourceConfig(typeName string, name string) (config string, err error)
}

type mqConsumerManager struct {
	handler   MQConsumerHandler
	services  map[string]spService
	consumers map[string]*mqConsumer
	mutex     sync.Mutex
	Log       logger.ILogger
}

func NewConsumerManager(handler MQConsumerHandler) (m *mqConsumerManager, err error) {
	m = &mqConsumerManager{handler: handler, consumers: make(map[string]*mqConsumer)}
	m.Log, err = logger.New("mq consumer", true)
	return
}

func (m *mqConsumerManager) Reset(services map[string]spService) (err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	//关闭已启动的服务
	for k, v := range m.consumers {
		if _, ok := services[k]; !ok {
			v.Stop()
			delete(m.consumers, k)
		}
	}

	//启动已添加的服务
	for k, v := range services {
		if _, ok := m.consumers[k]; !ok {
			paramMap, err := utility.GetParamsMap(v.Params)
			if err != nil {
				m.Log.Fatal("mq params error in:", v.Params, " error:", err)
				continue
			}
			mqParam := fmt.Sprintf("%s", paramMap["mq"])
			config, err := m.handler.GetSourceConfig("mq", mqParam)
			if err != nil || strings.EqualFold(config, "") {
				m.Log.Fatal("params must contain 'mq' in", v.Params)
				continue
			}
			m.consumers[k], err = NewMQConsumer(v.Params, config, func(h *msgHandler) bool {
				return m.handler.Hande(v.Script, v.Params, h)
			})
			if err != nil {
				m.Log.Fatal("mq create error:", err)
				continue
			}
			fmt.Println("::Start MQ Consumer:", mqParam)
			go m.consumers[k].Start()
		}
	}
	return
}

func (m *mqConsumerManager) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for k, v := range m.consumers {
		v.Stop()
		delete(m.consumers, k)
	}
}
