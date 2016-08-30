package main

import (
	"fmt"

	"github.com/arsgo/lib4go/db"
	"github.com/arsgo/lib4go/script"
	lua "github.com/yuin/gopher-lua"
)

//db操作类，用于lua脚本直接调用
//local cdb=db.new("agt_comm") ---根据zk中配置的数据连接名称创建DB实例
//local result,err=cdb:query("select 1 from dual where id=@id",{id=1})
//local result,err=cdb:execute("update tb1 set time=sysdate where id=@id",{id=1})
//local trans=cdb:begin()
//local result,err=trans:cdb:execute("update tb1 set time=sysdate where id=@id",{id=1})
//trans:rollback()
//trans.commit()

func (s *SPServer) getdbTypeBinder() []script.LuaTypesBinder {
	return []script.LuaTypesBinder{
		script.LuaTypesBinder{
			Name:    "dbtrans",
			NewFunc: map[string]lua.LGFunction{},
			Methods: map[string]lua.LGFunction{
				"execute":  typeDBTransExecute,
				"query":    typeDBTransQuery,
				"scalar":   typeDBTransScalar,
				"commit":   typeDBTransCommit,
				"rollback": typeDBTransRollback,
			},
		},
		script.LuaTypesBinder{
			Name: "db",
			NewFunc: map[string]lua.LGFunction{
				"new": s.typeDBType,
			},
			Methods: map[string]lua.LGFunction{
				"execute":   typeDBExecute,
				"begin":     typeDBBegin,
				"executeSP": typeDBExecuteSP,
				"query":     typeDBQuery,
				"scalar":    typeDBScalar,
			},
		}}
}

// Constructor
func (s *SPServer) typeDBType(L *lua.LState) int {
	var err error
	ud := L.NewUserData()
	name := L.CheckString(1)
	ud.Value, err = s.NewDB(name)
	if err != nil {
		return pushValues(L, "", err)
	}
	L.SetMetatable(ud, L.GetTypeMetatable("db"))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkDBType(L *lua.LState) *db.DBScriptBind {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*db.DBScriptBind); ok {
		return v
	}
	L.RaiseError("bad argument  (db expected, got %s)", ud.Type().String())
	return nil
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkDBTransType(L *lua.LState) *db.DBScriptBindTrans {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*db.DBScriptBindTrans); ok {
		return v
	}
	L.RaiseError("bad argument  (db trans expected, got %s)", ud.Type().String())
	return nil
}

func typeDBExecute(L *lua.LState) int {
	p := checkDBType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	a, b, c, d := p.Execute(query, input)
	return pushValues(L, a, b, c, d)
}
func typeDBExecuteSP(L *lua.LState) int {
	p := checkDBType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	a, b, c, d := p.ExecuteSP(query, input)
	return pushValues(L, a, b, c, d)
}
func typeDBQuery(L *lua.LState) int {
	p := checkDBType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	a, b, c, d := p.Query(query, input)
	return pushValues(L, a, b, c, d)
}
func typeDBScalar(L *lua.LState) int {
	p := checkDBType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	a, b, c, d := p.Scalar(query, input)
	return pushValues(L, a, b, c, d)
}
func typeDBBegin(L *lua.LState) int {
	p := checkDBType(L)
	ts, err := p.Begin()
	if err != nil {
		return pushValues(L, "", err)
	}
	ud := L.NewUserData()
	ud.Value = ts
	L.SetMetatable(ud, L.GetTypeMetatable("dbtrans"))
	L.Push(ud)
	return 1
}

func typeDBTransExecute(L *lua.LState) int {
	p := checkDBTransType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	a, b, c, d := p.Execute(query, input)
	return pushValues(L, a, b, c, d)
}

func typeDBTransQuery(L *lua.LState) int {
	p := checkDBTransType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	a, b, c, d := p.Query(query, input)
	return pushValues(L, a, b, c, d)
}
func typeDBTransScalar(L *lua.LState) int {
	p := checkDBTransType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	a, b, c, d := p.Scalar(query, input)
	return pushValues(L, a, b, c, d)
}
func typeDBTransCommit(L *lua.LState) int {
	p := checkDBTransType(L)
	a := p.Commit()
	return pushValues(L, a)
}
func typeDBTransRollback(L *lua.LState) int {
	p := checkDBTransType(L)
	a := p.Rollback()
	return pushValues(L, a)
}

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
