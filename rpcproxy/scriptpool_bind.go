package rpcproxy

import (
	"errors"

	"github.com/colinyl/lib4go/mem"
	"github.com/colinyl/lib4go/mq"
	"github.com/colinyl/lib4go/net"
	"github.com/colinyl/lib4go/utility"
)

//NewRPCClient 创建RPC client
func (s *ScriptPool) NewRPCClient() (*RPCBinder, error) {
	if s.rpcclient == nil {
		return nil, errors.New("not support rpc client")
	}
	return NewRPCBind(s.rpcclient), nil
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

//NewHTTPClient http client
func (s *ScriptPool) NewHTTPClient() *net.HTTPClient {
	return net.NewHTTPClient()
}

//bindGlobalLibs 绑定lib
func (s *ScriptPool) bindGlobalLibs(extlibs map[string]interface{}) (funs map[string]interface{}) {
	funs = map[string]interface{}{
		"print":         s.Log.Info,
		"printf":        s.Log.Infof,
		"error":         s.Log.Error,
		"errorf":        s.Log.Errorf,
		"NewGUID":       utility.GetGUID,
		"NewRPC":        s.NewRPCClient,
		"NewMQProducer": s.NewMQProducer,
		"NewMemcached":  s.NewMemcached,
		"NewXHttp":      s.NewHTTPClient,
		"NewSecurity":   s.NewBindSecurity,
	}
	for i, v := range extlibs {
		funs[i] = v
	}
	return
}
