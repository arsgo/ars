package rpcproxy

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/elastic"
	"github.com/colinyl/lib4go/influxdb"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/mem"
	"github.com/colinyl/lib4go/mq"
	"github.com/colinyl/lib4go/script"
	"github.com/colinyl/lib4go/utility"
)

//scriptInputArgs 脚本输入参数
type scriptInputArgs struct {
	Input  json.RawMessage `json:"input"`
	Params json.RawMessage `json:"params"`
}

//getScriptInputArgs 获取脚本输入参数
func getScriptInputArgs(input string, params string) string {
	args := scriptInputArgs{}
	args.Input = []byte(input)
	text, err := utility.GetParams(params)
	if err == nil {
		args.Params = []byte(text)
	}
	buffer, _ := json.Marshal(&args)
	return string(buffer)
}

//ScriptPool 创建ScriptPool
type ScriptPool struct {
	Pool          *script.LuaPool
	Log           *logger.Logger
	clusterClient cluster.IClusterClient
	rpcclient     *RPCClient
}

//NewScriptPool 创建脚本POOl
func NewScriptPool(clusterClient cluster.IClusterClient, rpcclient *RPCClient) (p *ScriptPool, err error) {
	p = &ScriptPool{}
	p.clusterClient = clusterClient
	p.rpcclient = rpcclient
	p.Pool = script.NewLuaPool()
	p.Pool.SetPackages(`E:\Project\golang\bin\scripts\xlib`, `E:\Project\golang\bin\scripts`)
	p.Log, err = logger.New("script", true)
	p.Pool.RegisterLibs(p.bindGlobalLibs())
	return
}

//Call 执行脚本
func (s *ScriptPool) Call(name string, input string, params string) ([]string, error) {
	if strings.EqualFold(name, "") {
		return nil, errors.New("script is nil")
	}
	script := name
	if !strings.HasPrefix(name, "./") {
		script = "./" + strings.TrimLeft(name, "/")
	}
	return s.Pool.Call(script, getScriptInputArgs(input, params))
}

//NewRPCClient 创建RPC client
func (s *ScriptPool) NewRPCClient() (*RPCClient, error) {
	if s.rpcclient == nil {
		return nil, errors.New("not support rpc client")
	}
	return s.rpcclient, nil
}

//NewInfluxDB 创建InfluxDB操作对象
func (s *ScriptPool) NewInfluxDB(name string) (p *influxdb.InfluxDB, err error) {
	config, err := s.clusterClient.GetDBConfig(name)
	if err != nil {
		return
	}
	p, err = influxdb.New(config)
	return
}

//NewMemcached 创建Memcached对象
func (s *ScriptPool) NewMemcached(name string) (p *mem.MemcacheClient, err error) {
	config, err := s.clusterClient.GetDBConfig(name)
	if err != nil {
		return
	}
	p, err = mem.New(config)
	return
}

//NewMQProducer 创建MQ Producer对象
func (s *ScriptPool) NewMQProducer(name string) (p *mq.MQProducer, err error) {
	config, err := s.clusterClient.GetMQConfig(name)
	if err != nil {
		return
	}
	p, err = mq.NewMQProducer(config)
	return
}

//NewElastic 创建Elastic对象
func (s *ScriptPool) NewElastic(name string) (es *elastic.ElasticSearch, err error) {
	config, err := s.clusterClient.GetElasticConfig(name)
	if err != nil {
		return
	}
	es, err = elastic.New(config)
	return
}

//bindGlobalLibs 绑定lib
func (s *ScriptPool) bindGlobalLibs() (funs map[string]interface{}) {
	funs = map[string]interface{}{
		"print":         s.Log.Info,
		"printf":        s.Log.Infof,
		"error":         s.Log.Error,
		"errorf":        s.Log.Errorf,
		"NewRPC":        s.NewRPCClient,
		"NewMQProducer": s.NewMQProducer,
		"NewElastic":    s.NewElastic,
		"NewInfluxDB":   s.NewInfluxDB,
		"NewMemcached":  s.NewMemcached,
	}
	return
}
