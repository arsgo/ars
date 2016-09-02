package script

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/ars/rpc"
	"github.com/arsgo/lib4go/concurrent"
	"github.com/arsgo/lib4go/logger"
	spt "github.com/arsgo/lib4go/script"
	"github.com/arsgo/lib4go/utility"
)

//scriptInputArgs 脚本输入参数
type scriptInputArgs struct {
	Input  json.RawMessage `json:"input"`
	Params json.RawMessage `json:"params"`
}

//getScriptInputArgs 获取脚本输入参数
func getScriptInputArgs(input string, params string) (r string) {
	if strings.EqualFold(input, "") {
		input = "{}"
	}
	args := scriptInputArgs{}
	args.Input = []byte(input)
	text, err := utility.GetParams(params)
	if strings.EqualFold(text, "") {
		text = "{}"
	}
	args.Params = []byte(text)
	buffer, err := json.Marshal(&args)
	if err != nil {
		fmt.Printf("get script args error:%v,%v\n", err, args)
	}
	r = string(buffer)
	return
}

//ScriptPool 创建ScriptPool
type ScriptPool struct {
	pool          *spt.LuaPool
	Log           logger.ILogger
	clusterClient cluster.IClusterClient
	rpcclient     *rpc.RPCClient
	snaps         *concurrent.ConcurrentMap
	mqservices    *concurrent.ConcurrentMap
	collectors    map[string]base.ICollector
}

//NewScriptPool 创建脚本POOl
func NewScriptPool(clusterClient cluster.IClusterClient, rpcclient *rpc.RPCClient, types []spt.LuaTypesBinder,
	loggerName string, collectors map[string]base.ICollector) (p *ScriptPool, err error) {
	p = &ScriptPool{snaps: concurrent.NewConcurrentMap(), collectors: collectors}
	p.mqservices = concurrent.NewConcurrentMap()
	p.clusterClient = clusterClient
	p.rpcclient = rpcclient
	p.pool = spt.NewLuaPool()
	p.Log, err = logger.Get(loggerName)
	p.pool.Binder.RegisterModules(p.bindModules())
	p.pool.Binder.RegisterGlobal(p.bindGlobal())
	types = append(types, p.bindTypes()...)
	p.pool.Binder.RegisterTypes(types...)
	return
}
func (s *ScriptPool) SetPoolSize(minSize int, maxSize int) {
	s.pool.SetPoolSize(minSize, maxSize)
}
func (s *ScriptPool) PreLoad(script string, minSize int, maxSize int) error {
	return s.pool.PreLoad(utility.GetExcPath(script, "bin"), minSize, maxSize)
}

func (s *ScriptPool) SetPackages(path ...string) {
	if len(path) == 0 {
		return
	}
	npath := make([]string, 0, len(path))
	for _, v := range path {
		npath = append(npath, utility.GetExcPath(v, "bin"))
	}
	s.pool.Binder.SetPackages(npath...)
}
func (s *ScriptPool) createSnap(p ...interface{}) (interface{}, error) {
	ss := &base.ProxySnap{}
	ss.ElapsedTime = base.ServerSnap{}
	return ss, nil
}
func (s *ScriptPool) setLifeTime(name string, start time.Time) {
	_, snap, _ := s.snaps.GetOrAdd(name, s.createSnap)
	if snap == nil {
		return
	}
	snap.(*base.ProxySnap).ElapsedTime.Add(start)
}

//Call 执行脚本
func (s *ScriptPool) Call(name string, context base.InvokeContext) ([]string, map[string]string, error) {
	defer base.RunTime("script poll call", time.Now())
	if strings.EqualFold(name, "") {
		return nil, nil, errors.New("script is nil")
	}
	script := utility.GetExcPath(name, "bin")
	defer s.setLifeTime(script, time.Now())
	input := spt.InputArgs{Script: script, Session: context.Session, Body: context.Body, TaskType: context.TaskType, TaskName: context.TaskName}
	input.Input = getScriptInputArgs(context.Input, context.Params)
	return s.pool.Call(input)
}

//GetSnap 获取当前脚本
func (s *ScriptPool) GetSnap() (r []interface{}) {
	poolSnaps := s.pool.GetSnap()
	snaps := s.snaps.GetAll()
	return base.GetProxySnap(poolSnaps, snaps)
}

//Close 关闭脚本引擎
func (s *ScriptPool) Close() {
	s.pool.Close()
}
