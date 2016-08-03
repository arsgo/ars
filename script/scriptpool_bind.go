package script

import (
	"errors"
	"time"

	"github.com/arsgo/lib4go/mem"
	m "github.com/arsgo/lib4go/mq"
	"github.com/arsgo/lib4go/net"
	"github.com/arsgo/lib4go/security/weixin"
	"github.com/arsgo/lib4go/utility"
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
func (s *ScriptPool) NewMQProducer(name string) (p *m.MQProducer, err error) {
	config, err := s.clusterClient.GetMQConfig(name)
	if err != nil {
		return
	}
	pd, err := s.mqservices.GetOrAdd(name, func(p ...interface{}) (interface{}, error) {
		return m.NewMQProducer(p[0].(string))
	}, config)
	p = pd.(*m.MQProducer)
	return
}

//NewHTTPClient http client
func (s *ScriptPool) NewHTTPClient() *net.HTTPClient {
	return net.NewHTTPClient()
}

//NewHTTPClientCert http client
func (s *ScriptPool) NewHTTPClientCert(certFile string, keyFile string, caFile string) (*net.HTTPClient, error) {
	return net.NewHTTPClientCert(certFile, keyFile, caFile)
}

//NewHTTPClientProxy 根据代理服务器地址创建http client,代理服务器格式:http://192.168.101.1:8080
func (s *ScriptPool) NewHTTPClientProxy(proxy string) *net.HTTPClient {
	return net.NewHTTPClientProxy(proxy)
}

//NewWechat 创建微信加解密对象
func (s *ScriptPool) NewWechat(appid string, token string, encodingAESKey string) (weixin.Wechat, error) {
	return weixin.NewWechat(appid, token, encodingAESKey)
}

//Sleep 休息指定时间
func (s *ScriptPool) Sleep(r int) {
	time.Sleep(time.Second * time.Duration(r))
}

//bindGlobalLibs 绑定lib
func (s *ScriptPool) bindGlobalLibs(extlibs map[string]interface{}) (funs map[string]interface{}) {
	funs = map[string]interface{}{
		"NewGUID":            utility.GetGUID,
		"NewRPC":             s.NewRPCClient,
		"NewMQProducer":      s.NewMQProducer,
		"NewMemcached":       s.NewMemcached,
		"NewXHttp":           s.NewHTTPClient,
		"NewHTTPClientCert":  s.NewHTTPClientCert,
		"NewHTTPClientProxy": s.NewHTTPClientProxy,
		"NewSecurity":        s.NewBindSecurity,
		"NewWechat":          s.NewWechat,
	}
	for i, v := range extlibs {
		funs[i] = v
	}
	return
}
