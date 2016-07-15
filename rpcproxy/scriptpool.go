package rpcproxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/colinyl/ars/base"
	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/concurrent"
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
	Log           logger.ILogger
	clusterClient cluster.IClusterClient
	rpcclient     *RPCClient
	snaps         concurrent.ConcurrentMap
}

//NewScriptPool 创建脚本POOl
func NewScriptPool(clusterClient cluster.IClusterClient, rpcclient *RPCClient, extlibs map[string]interface{},
	loggerName string) (p *ScriptPool, err error) {
	p = &ScriptPool{snaps: concurrent.NewConcurrentMap()}
	p.clusterClient = clusterClient
	p.rpcclient = rpcclient
	p.Pool = script.NewLuaPool()
	p.Pool.SetPackages(`./scripts`, `./script`, `./scripts/xlib`,`./scripts/lib`, `./scripts/lib/xlib`, `./script/lib`, `./script/lib/xlib`)
	p.Log, err = logger.Get(loggerName, true)
	p.Pool.RegisterLibs(p.bindGlobalLibs(extlibs))
	p.Pool.RegisterModules(p.bindModules())
	return
}
func (s *ScriptPool) setLifeTime(name string, start time.Time) {
	ss := &ProxySnap{}
	ss.ElapsedTime = ServerSnap{}
	snap := s.snaps.GetOrAdd(name, ss)
	snap.(*ProxySnap).ElapsedTime.Add(start)
}

//Call 执行脚本
func (s *ScriptPool) Call(name string, context base.InvokeContext) ([]string, map[string]string, error) {
	if strings.EqualFold(name, "") {
		return nil, nil, errors.New("script is nil")
	}
	script := name
	//if !strings.HasPrefix(name, "./") {
	//script = "./" + strings.TrimLeft(name, "/")
	//}
	defer s.setLifeTime(script, time.Now())
	return s.Pool.Call(script, context.Session, getScriptInputArgs(context.Input, context.Params), context.Body)
}

//GetSnap 获取当前脚本
func (s *ScriptPool) GetSnap() (r []interface{}) {
	poolSnaps := s.Pool.GetSnap()
	snaps := s.snaps.GetAll()
	return getProxySnap(poolSnaps, snaps)
}

//Close 关闭脚本引擎
func (s *ScriptPool) Close() {
	s.Pool.Close()
}
