package rpcproxy

import (
	"github.com/colinyl/lib4go/utility"
	"github.com/yuin/gopher-lua"
)

func (s *ScriptPool) moduleGetGUID(ls *lua.LState) int {
	return pushValues(ls, utility.GetGUID())
}
