package mqservice

import (
	"github.com/colinyl/ars/scriptservice"
	lp "github.com/colinyl/lib4go/lua"
	"github.com/colinyl/stomp"
	"github.com/yuin/gopher-lua"
)

type ConfigHandler interface {
	GetMQConfig(string) (string, error)
}
type MQBinder struct {
	handler ConfigHandler
	pool    *lp.LuaPool
}

func NewMQBinder(c ConfigHandler, pool *lp.LuaPool) *MQBinder {
	return &MQBinder{handler: c, pool: pool}
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
				p.Consume(L.CheckString(2), func(msg stomp.MsgHandler)bool {
					c.pool.Call(L.CheckString(3), msg.GetMessage())
					return true
				})
				return result
			},
		}})
}
