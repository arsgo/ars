package rpcproxy

import "github.com/yuin/gopher-lua"

//Request RPC Reuqest调用
func (s *ScriptPool) Request(L *lua.LState) int {
    if L.GetTop()!=2{
        L.Push(lua.LString(""))
        L.Push(lua.LString("输入参数个数有误"))
        return 2
    }
	cmd := L.CheckString(1)
	input := L.CheckString(2)
	r, e := s.rpcclient.Request(cmd, input)
	L.Push(lua.LString(r))
	if e != nil {
		L.Push(lua.LString(e.Error()))
		return 2
	}
    return 1
}
