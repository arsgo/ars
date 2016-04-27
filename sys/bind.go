package sys

import l "github.com/yuin/gopher-lua"

func SysInfoLoader(L *l.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]l.LGFunction{
	"cpu":  getCPU,
	"mem":  getMemory,
	"disk": getDisk,
}

func getCPU(L *l.LState) int {
	L.Push(l.LString(GetCPU()))
	return 1
}
func getMemory(L *l.LState) int {
	L.Push(l.LString(GetMemory()))
	return 1
}
func getDisk(L *l.LState) int {
	L.Push(l.LString(GetDisk()))
	return 1
}
