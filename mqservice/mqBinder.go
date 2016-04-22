package mqservice

import (
	"fmt"

	"github.com/colinyl/ars/scriptservice"
	"github.com/colinyl/stomp"
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
func (c *MQBinder) BindMQService(L *lua.LState) {
	scriptservice.Bind(L, &scriptservice.ScriptBindClass{ClassName: "mq",
		ConstructorName: "new",
		ConstructorFunc: func(L *lua.LState) interface{} {
			config, _ := c.handler.GetMQConfig(L.CheckString(1))
			return NewMQService(config)
		}, ObjectMethods: map[string]scriptservice.ScriptBindFunc{
			"close": func(L *lua.LState) (result []string) {
				if L.GetTop() != 1 {
					result = append(result, "input args error")
					return
				}
				ud := L.CheckUserData(1)
				if _, ok := ud.Value.(IMQService); !ok {
					result = append(result, "MQService expected")
					return
				}
				p := ud.Value.(IMQService)
				p.Close()
				return
			},
			"send": func(L *lua.LState) (result []string) {
				if L.GetTop() != 3 {
					result = append(result, "input args error")
					return
				}
				ud := L.CheckUserData(1)
				if _, ok := ud.Value.(IMQService); !ok {
					result = append(result, "MQService expected")
					return
				}
				p := ud.Value.(IMQService)
				err := p.Send(L.CheckString(2), L.CheckString(3))
				if err != nil {
					result = append(result, err.Error())
				}
				return result
			},
			"consume": func(L *lua.LState) (result []string) {
				if L.GetTop() != 3 {
					result = append(result, "input args error")
					return
				}
				ud := L.CheckUserData(1)
				if _, ok := ud.Value.(IMQService); !ok {
					result = append(result, "MQConsumer expected")
					return
				}
				p := ud.Value.(IMQService)
				fmt.Println("call func")
				p.Consume(L.CheckString(2), func(msg stomp.MsgHandler) {
					L.CallByParam(lua.P{
						Fn:      L.CheckFunction(3),
						NRet:    0,
						Protect: true},
						lua.LString(msg.GetMessage()))
					msg.Ack()
				})
				return result
			},
		}})
}
