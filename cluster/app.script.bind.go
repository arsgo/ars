package cluster

import (
	"errors"
	"strings"

	"github.com/colinyl/lib4go/script"
)

type scriptEngine struct {
	pool *script.LuaPool
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

func (a *appServer) bindGlobalLibs() (funs map[string]interface{}) {
	funs = map[string]interface{}{
		"print":  a.Log.Info,
		"printf": a.Log.Infof,
		"error":  a.Log.Error,
		"errorf": a.Log.Errorf,
		"NewRPC": a.NewRpcHandler,
	}
	return
}

func NewScriptEngine(app *appServer) *scriptEngine {
	pool := script.NewLuaPool()
	pool.RegisterLibs(app.bindGlobalLibs())
	return &scriptEngine{pool: pool}
}
