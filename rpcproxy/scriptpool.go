package rpcproxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/script"
	"github.com/colinyl/lib4go/utility"
)

//scriptInputArgs 脚本输入参数
type scriptInputArgs struct {
	Input  json.RawMessage `json:"input"`
	Params json.RawMessage `json:"params"`
}

//getScriptInputArgs 获取脚本输入参数
func getScriptInputArgs(input string, params string) (r string) {
	args := scriptInputArgs{}
	args.Input = []byte(input)
	text, err := utility.GetParams(params)
	if strings.EqualFold(text, "") {
		text = "{}"
	}
	args.Params = []byte(text)

	buffer, err := json.Marshal(&args)
	if err != nil {
		fmt.Println(err)
	}
	r = string(buffer)
	return
}

//ScriptPool 创建ScriptPool
type ScriptPool struct {
	Pool          *script.LuaPool
	Log           *logger.Logger
	clusterClient cluster.IClusterClient
	rpcclient     *RPCClient
}

//NewScriptPool 创建脚本POOl
func NewScriptPool(clusterClient cluster.IClusterClient, rpcclient *RPCClient, extlibs map[string]interface{}) (p *ScriptPool, err error) {
	p = &ScriptPool{}
	p.clusterClient = clusterClient
	p.rpcclient = rpcclient
	p.Pool = script.NewLuaPool()
	p.Pool.SetPackages(`./scripts/xlib`, `./scripts`)
	p.Log, err = logger.New("script", true)
	p.Pool.RegisterLibs(p.bindGlobalLibs(extlibs))
	return
}

//Call 执行脚本
func (s *ScriptPool) Call(name string, input string, params string) ([]string, map[string]string, error) {
	if strings.EqualFold(name, "") {
		return nil, nil, errors.New("script is nil")
	}
	script := name
	if !strings.HasPrefix(name, "./") {
		script = "./" + strings.TrimLeft(name, "/")
	}
	return s.Pool.Call(script, getScriptInputArgs(input, params))
}
