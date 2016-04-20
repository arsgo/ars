package scriptservice

import "github.com/yuin/gopher-lua"

type ScriptBindFunc func(*lua.LState) []string
type ScriptBinderFuncs map[string]ScriptBindFunc

type ScriptBindClass struct {
	ClassName     string
	ObjectFuncs   map[string]func(*lua.LState) interface{}
	Funcs         ScriptBinderFuncs
	ObjectMethods ScriptBinderFuncs
}

func Bind(L *lua.LState, pk *ScriptBindClass) {
	mt := L.NewTypeMetatable(pk.ClassName)
	L.SetGlobal(pk.ClassName, mt)
	for oName, oFunc := range pk.ObjectFuncs {
		L.SetField(mt, oName, L.NewFunction(func(ls *lua.LState) int {
			ud := ls.NewUserData()
			ud.Value = oFunc(ls)
			ls.SetMetatable(ud, ls.GetTypeMetatable(pk.ClassName))
			ls.Push(ud)
			return 1
		}))
	}
	for cName, cFunc := range pk.Funcs {
		L.SetField(mt, cName, L.NewFunction(getFunc(L, cFunc)))
	}
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), getFuncMap(L, pk.ObjectMethods)))
}
func getFunc(L *lua.LState, fun ScriptBindFunc) lua.LGFunction {
	return func(ls *lua.LState) int {
		results := fun(ls)
		for _, v := range results {
			ls.Push(lua.LString(v))
		}
		return len(results)
	}
}
func getFuncMap(L *lua.LState, funs ScriptBinderFuncs) (rfun map[string]lua.LGFunction) {
	rfun = make(map[string]lua.LGFunction)
	for i, v := range funs {
		rfun[i] = getFunc(L, v)
	}
	return
}
