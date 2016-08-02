package script

import (
	"encoding/json"

	"github.com/arsgo/ars/rpc"
	"github.com/arsgo/lib4go/utility"
	"github.com/yuin/gopher-lua"
)

type RPCBinder struct {
	client *rpc.RPCClient
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

func NewRPCBind(client *rpc.RPCClient) *RPCBinder {
	return &RPCBinder{client: client}
}

func (b *RPCBinder) AsyncRequest(name string, tb *lua.LTable) (s string, err error) {
	input, err := luaTable2Json(tb)
	if err != nil {
		return
	}
	return b.client.AsyncRequest(name, input, utility.GetSessionID())
}

func (b *RPCBinder) GetAsyncResult(session string) (s interface{}, err interface{}) {
	return b.client.GetAsyncResult(session)
}
func (b *RPCBinder) Request(name string, tb *lua.LTable) (s string, err error) {
	input, err := luaTable2Json(tb)
	if err != nil {
		return
	}
	return b.client.Request(name, input, utility.GetSessionID())
}
