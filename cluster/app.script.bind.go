package cluster

import (
	"errors"
	"strings"

	"github.com/colinyl/ars/scriptservice"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/lua"
	"github.com/colinyl/lib4go/mq"
	"github.com/colinyl/lib4go/sysinfo"
	l "github.com/yuin/gopher-lua"
)

type scriptEngine struct {
	pool *lua.LuaPool
}
type scriptCallbackHandler struct {
	queue chan []string
	Log   *logger.Logger
}

func (h *scriptCallbackHandler) GetResult() (result []string) {
	result = <-h.queue
	return
}

func (e *scriptEngine) Call(name string, input string) ([]string, error) {
	if strings.EqualFold(name, "") {
		return nil, errors.New("script is nil")
	}
	script := strings.Replace(name, ".", "/", -1)
	script = strings.Replace(script, "\\", "/", -1)
	if !strings.HasPrefix(script, "./") {
		script = "./" + strings.TrimLeft(name, "/")
	}
	return e.pool.Call(script, input)
}

func handerGet(ls *l.LState) (result []string) {
	ud := ls.CheckUserData(1)
	if _, ok := ud.Value.(*scriptCallbackHandler); !ok {
		result = append(result, "rpc handler expected")
		return
	}
	p := ud.Value.(*scriptCallbackHandler)
	result = p.GetResult()
	return result
}
func (a *appServer) bindRPCRequestService(L *l.LState) {
	scriptservice.Bind(L, &scriptservice.ScriptBindClass{ClassName: "rpcRequest",
		ConstructorName: "async",
		ConstructorFunc: func(ls *l.LState) interface{} {
			return a.asyncRequest(ls)
		}, ObjectMethods: map[string]scriptservice.ScriptBindFunc{
			"get": handerGet,
		}})
}
func (a *appServer) bindRPCSendService(L *l.LState) {
	scriptservice.Bind(L, &scriptservice.ScriptBindClass{ClassName: "rpcSend",
		ConstructorName: "async",
		ConstructorFunc: func(L *l.LState) interface{} {
			return a.asyncSend(L)
		}, ObjectMethods: map[string]scriptservice.ScriptBindFunc{
			"get": handerGet,
		}})
}
func (a *appServer) bindLogger() (fn []lua.Luafunc) {
	fn = append(fn, lua.Luafunc{
		Name: "print",
		Function: func(L *l.LState) int {
			msg := L.CheckString(1)
			a.Log.Info(msg)
			return 0
		},
	})
	fn = append(fn, lua.Luafunc{
		Name: "error",
		Function: func(L *l.LState) int {
			msg := L.CheckString(1)
			a.Log.Error(msg)
			return 0
		},
	})
	fn = append(fn, lua.Luafunc{
		Name:     "info",
		Function: sysinfo.SysInfoLoader,
	})

	return
}

func NewScriptEngine(app *appServer) *scriptEngine {
	pool := lua.NewLuaPool(app.bindLogger()...)
	pool.AddUserData(app.bindRPCRequestService)
	pool.AddUserData(app.bindRPCSendService)
	mqBinder := mq.NewMQBinder(app.zkClient, pool)
	pool.AddUserData(mqBinder.BindMQService)
	return &scriptEngine{pool: pool}
}
