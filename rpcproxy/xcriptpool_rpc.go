package rpcproxy

import (
	"encoding/json"

	"github.com/yuin/gopher-lua"
)

type RPCBinder struct {
	client *RPCClient
}

func luaTable2Json(tb *lua.LTable) (s string, err error) {
	data := make(map[string]interface{})
	tb.ForEach(func(key lua.LValue, value lua.LValue) {
		data[key.String()] = value.String()
	})
	buffer, err := json.Marshal(&data)
	if err != nil {
		return
	}
	s = string(buffer)
	return
}

func NewRPCBind(client *RPCClient) *RPCBinder {
	return &RPCBinder{client: client}
}

func (b *RPCBinder) AsyncRequest(name string, tb *lua.LTable) (s string, err error) {
	input, err := luaTable2Json(tb)
	if err != nil {
		return
	}
	return b.client.AsyncRequest(name, input)
}
func (b *RPCBinder) Request(name string, tb *lua.LTable) (s string, err error) {
	input, err := luaTable2Json(tb)
	if err != nil {
		return
	}
	return b.client.Request(name, input)
}
