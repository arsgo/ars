package mqservice

import (
	"github.com/colinyl/ars/scriptservice"
	"github.com/yuin/gopher-lua"
)

type ConfigHandler interface {
	GetMQConfig(string) (string, error)
}
type MQBinder struct {
	handler ConfigHandler
}

func NewMQBinder(c ConfigHandler) *MQBinder {
	return &MQBinder{handler: c}
}
func (c *MQBinder) BindMQConsumerService(L *lua.LState) {
	scriptservice.Bind(L, &scriptservice.ScriptBindClass{ClassName: "mqc",
		ConstructorName: "new",
		ConstructorFunc: func(L *lua.LState) interface{} {
			config, _ := c.handler.GetMQConfig(L.CheckString(1))
			return NewMQService(config).NewConsumer()
		}, ObjectMethods: map[string]scriptservice.ScriptBindFunc{
			"consume": func(L *lua.LState) (result []string) {
				if L.GetTop() != 2 {
					result = append(result, "input args error")
					return
				}
				ud := L.CheckUserData(1)
				if _, ok := ud.Value.(IMQConsumer); !ok {
					result = append(result, "MQConsumer expected")
					return
				}
				p := ud.Value.(IMQConsumer)
				p.Consume(func(msg string) {
					L.CallByParam(lua.P{
						Fn:      L.CheckFunction(2),
						NRet:    0,
						Protect: true},
						lua.LString(msg))
				})
				return result
			},
		}})
}

func (c *MQBinder) BindMQPublisherService(L *lua.LState) {
	scriptservice.Bind(L, &scriptservice.ScriptBindClass{ClassName: "mqp",
		ConstructorName: "new",
		ConstructorFunc: func(L *lua.LState) interface{} {
			config, _ := c.handler.GetMQConfig(L.CheckString(1))
			return NewMQService(config).NewPublisher()
		}, ObjectMethods: map[string]scriptservice.ScriptBindFunc{
			"publish": func(L *lua.LState) (result []string) {
				if L.GetTop() != 2 {
					result = append(result, "input args error")
					return
				}
				ud := L.CheckUserData(1)
				if _, ok := ud.Value.(IMQPublisher); !ok {
					result = append(result, "MQPublisher expected")
					return result
				}
				p := ud.Value.(IMQPublisher)
				err := p.Publish(L.CheckString(2))
				if err != nil {
					result = append(result, err.Error())
				}
				return result
			},
		}})
}
