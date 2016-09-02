package script

import (
	m "github.com/arsgo/lib4go/mq"
	"github.com/arsgo/lib4go/script"
	lua "github.com/yuin/gopher-lua"
)

//mq操作类，用于lua脚本直接调用
//local producer=mq.new("mq")  --根据zk配置名称初始化mq
//producer:send(queue,content,timeout) --发送消息

func (s *ScriptPool) getMQTypeBinder() script.LuaTypesBinder {
	return script.LuaTypesBinder{
		Name: "mq",
		NewFunc: map[string]lua.LGFunction{
			"new": s.typeMQType,
		},
		Methods: map[string]lua.LGFunction{
			"send": typeMQSend,
		},
	}
}

// Constructor
func (s *ScriptPool) typeMQType(L *lua.LState) int {
	var err error
	ud := L.NewUserData()
	name := L.CheckString(1)
	ud.Value, err = s.NewMQProducer(name)
	if err != nil {
		return pushValues(L, "", err)
	}
	L.SetMetatable(ud, L.GetTypeMetatable("mq"))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkMQType(L *lua.LState) *m.MQProducer {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*m.MQProducer); ok {
		return v
	}
	L.RaiseError("bad argument  (MQProducer expected, got %s)", ud.Type().String())
	return nil
}

func typeMQSend(L *lua.LState) int {
	p := checkMQType(L)
	queue := L.CheckString(2)
	content := L.CheckString(3)
	timeout := L.CheckInt(4)
	a := p.Send(queue, content, timeout)
	return pushValues(L, a)
}
