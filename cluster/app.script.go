package cluster

import (
	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/lua"
	l "github.com/yuin/gopher-lua"
)

type scriptEngine struct {
	pool *lua.LuaPool
}

func (a *appServer) rpcBind(L *l.LState) int {

	var exports = map[string]l.LGFunction{
		"request": a.request,
		"send":    a.send,
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

func (a *appServer) request(L *l.LState) int {
	name := L.ToString(1)
	input := L.ToString(2)
	result, err := a.rcServerPool.Request(a.rcServicesMap.Next("-"), name, input)
	L.Push(l.LString(result))
	if err != nil {
		L.Push(l.LString(err.Error()))
	} else {
		L.Push(l.LNil)
	}
	return 2
}
func (a *appServer) send(L *l.LState) int {
	name := L.ToString(1)
	input := L.ToString(2)
	buffer := []byte(L.ToString(3))
	group := a.rcServicesMap.Next("-")
	result, err := a.rcServerPool.Send(group, name, input, buffer)
	L.Push(l.LString(result))
	if err != nil {
		L.Push(l.LString(err.Error()))
	} else {
		L.Push(l.LNil)
	}
	return 2
}

var callback rpcservice.ServiceCallBack

func NewScriptEngine(app *appServer) *scriptEngine {
	rpcFunc := lua.Luafunc{Name: "rpc", Function: app.rpcBind}
	return &scriptEngine{pool: lua.NewLuaPool(rpcFunc)}
}
