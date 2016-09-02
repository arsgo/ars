package script

import (
	"github.com/arsgo/lib4go/script"
	"github.com/yuin/gopher-lua"
)

func pushValues(ls *lua.LState, values ...interface{}) int {
	for _, v := range values {
		ls.Push(script.New(ls, v))
	}
	return len(values)
}

//Request RPC Reuqest调用
func (s *ScriptPool) moduleRPCGetResult(ls *lua.LState) int {
	asyncSessionID := ls.CheckString(1)
	timeout := 5000
	if ls.GetTop() == 2 {
		timeout = ls.CheckInt(2)
	}
	r, e := s.rpcclient.GetAsyncResult(asyncSessionID, timeout)
	return pushValues(ls, r, e)
}

//Request RPC Reuqest调用
func (s *ScriptPool) moduleRPCAsyncRequest(ls *lua.LState) int {
	session := ls.GetGlobal("__session__").String()
	cmd := ls.CheckString(1)
	input := ls.CheckTable(2)
	inputString, err := luaTable2Json(input)
	if err != nil {
		return pushValues(ls, "", "输入参数错误")
	}
	r, e := s.rpcclient.AsyncRequest(cmd, inputString, session)
	if e != nil {
		return pushValues(ls, r, e.Error())
	}
	return pushValues(ls, r)
}

//Request RPC Reuqest调用
func (s *ScriptPool) moduleRPCRequest(ls *lua.LState) int {
	session := ls.GetGlobal("__session__").String()
	cmd := ls.CheckString(1)
	input := ls.CheckTable(2)
	inputString, err := luaTable2Json(input)
	if err != nil {
		return pushValues(ls, "", "输入参数错误")
	}
	r, e := s.rpcclient.Request(cmd, inputString, session)
	if e != nil {
		return pushValues(ls, r, e.Error())
	}
	return pushValues(ls, r)
}
