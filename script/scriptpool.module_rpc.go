package script

import (
	"fmt"

	"github.com/yuin/gopher-lua"
)

func pushValues(ls *lua.LState, values ...interface{}) int {
	for _, v := range values {
		if v != nil {
			ls.Push(lua.LString(fmt.Sprintf("%v", v)))
		} else {
			ls.Push(lua.LNil)
		}
	}
	return len(values)
}

//Request RPC Reuqest调用
func (s *ScriptPool) moduleRPCRequest(ls *lua.LState) int {
	session := ls.GetGlobal("__session").String()
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
