package cluster

import (
	"errors"
	"strings"

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
func (a *appServer) asyncRequest(L *l.LState) (handler *scriptCallbackHandler) {
	//	fmt.Println("--asyncRequest")
	a.Log.Info("--asyncRequest")
	name := L.ToString(1)
	input := L.ToString(2)
	handler = &scriptCallbackHandler{queue: make(chan []string, 1), Log: a.Log}
	go func(name string, input string) {
		handler.queue <- a.request(name, input)
	}(name, input)
	return
}
func (a *appServer) asyncSend(L *l.LState) (handler *scriptCallbackHandler) {
	handler = &scriptCallbackHandler{queue: make(chan []string, 1), Log: a.Log}
	go func() {
		handler.queue <- a.send(L)
	}()
	return
}

func (a *appServer) request(name string, input string) (result []string) {
	rest, err := a.rcServerPool.Request(a.rcServicesMap.Next("-"), name, input)
	result = append(result, rest)
	if err != nil {	
		result = append(result, err.Error())
	} else {
		result = append(result, "")
	}
	return
}

func (a *appServer) send(L *l.LState) (result []string) {
	name := L.ToString(1)
	input := L.ToString(2)
	buffer := []byte(L.ToString(3))
	group := a.rcServicesMap.Next("-")
	rest, err := a.rcServerPool.Send(group, name, input, buffer)
	result = append(result, rest)
	if err != nil {
		result = append(result, err.Error())
	} else {
		result = append(result, "")
	}
	return
}
func (a *appServer) bindRPCService(L *l.LState) {
	scriptservice.Bind(L, &scriptservice.ScriptBindClass{ClassName: "rpcfactory",
		ObjectFuncs: map[string]func(*l.LState) interface{}{
			"asyncRequest": func(ls *l.LState) interface{} {
				return a.asyncRequest(ls)
			},
			//"asyncSend": func(L *l.LState) interface{} {
			//	return a.asyncSend(L)
			//},
		}, ObjectMethods: map[string]scriptservice.ScriptBindFunc{
			"getResult": func(ls *l.LState) (result []string) {
				ud := ls.CheckUserData(1)
				if _, ok := ud.Value.(*scriptCallbackHandler); !ok {
					result = append(result, "rpc handler expected")
					return
				}
				p := ud.Value.(*scriptCallbackHandler)
				result = p.GetResult()
				return result
			},
		}})
}

func NewScriptEngine(app *appServer) *scriptEngine {
	pool := lua.NewLuaPool()
	pool.AddUserData(app.bindRPCService)
	//	pool.AddUserData(mqservice.BindMQService)
	return &scriptEngine{pool: pool}
}
