package script

import (
	"github.com/arsgo/lib4go/elastic"
	"github.com/arsgo/lib4go/script"
	lua "github.com/yuin/gopher-lua"
)

//elastic操作类，用于lua脚本直接调用
//local es,err=elastic.new("es")
//if err~=nil then
//	 print(err)
//end
func (s *ScriptPool) getElasticTypeBinder() script.LuaTypesBinder {
	return script.LuaTypesBinder{
		Name: "elastic",
		NewFunc: map[string]lua.LGFunction{
			"new": s.typeElasticType,
		},
		Methods: map[string]lua.LGFunction{
			"create": typeElasticCreate,
			"search": typeElasticSearch,
		},
	}
}

// Constructor
func (s *ScriptPool) typeElasticType(L *lua.LState) int {
	var err error
	ud := L.NewUserData()
	name := L.CheckString(1)
	ud.Value, err = s.NewElastic(name)
	if err != nil {
		return pushValues(L, "", err)
	}
	L.SetMetatable(ud, L.GetTypeMetatable("elastic"))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkElasticType(L *lua.LState) *elastic.ElasticSearch {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*elastic.ElasticSearch); ok {
		return v
	}
	L.RaiseError("bad argument  (elastic.ElasticSearch expected, got %s)", ud.Type().String())
	return nil
}
func typeElasticCreate(L *lua.LState) int {
	p := checkElasticType(L)
	name := L.CheckString(2)
	tp := L.CheckString(3)
	data := L.CheckString(4)
	a, err := p.Create(name, tp, data)
	return pushValues(L, a, err)
}
func typeElasticSearch(L *lua.LState) int {
	p := checkElasticType(L)
	name := L.CheckString(2)
	tp := L.CheckString(3)
	data := L.CheckString(4)
	a, err := p.Search(name, tp, data)
	return pushValues(L, a, err)
}
