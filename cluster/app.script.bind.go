package cluster

import (
	"errors"
	"strings"

	"github.com/colinyl/ars/mqservice"
	"github.com/colinyl/ars/scriptservice"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/lua"
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

func NewScriptEngine(app *appServer) *scriptEngine {
	pool := lua.NewLuaPool()
	pool.AddUserData(app.bindRPCRequestService)
	pool.AddUserData(app.bindRPCSendService)
	mqBinder := mqservice.NewMQBinder(app.zkClient)
	pool.AddUserData(mqBinder.BindMQConsumerService)
	pool.AddUserData(mqBinder.BindMQPublisherService)
	return &scriptEngine{pool: pool}
}
