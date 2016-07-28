package rpcservice

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	"github.com/colinyl/ars/servers/config"
	"github.com/colinyl/lib4go/concurrent"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/pool"
)

type rpcServerService struct {
	Status bool
	IP     string
}

type RPCServerPool struct {
	pool       *pool.ObjectPool
	servers    concurrent.ConcurrentMap
	Log        logger.ILogger
	loggerName string
	domain     string
	MaxRetry   int
	MinSize    int
	MaxSize    int
}

//NewRPCServerPool 创建RPC连接池
func NewRPCServerPool(minSize int, maxSize int, loggerName string) *RPCServerPool {
	var err error
	conf, _ := config.Get()
	pl := &RPCServerPool{MinSize: minSize, MaxSize: maxSize, MaxRetry: 1, loggerName: loggerName, domain: conf.Domain}
	pl.pool = pool.New()
	pl.servers = concurrent.NewConcurrentMap()
	pl.Log, err = logger.Get(loggerName)
	if err != nil {
		log.Println(err)
	}
	return pl
}

//ResetAllPoolSize 重置所有连接池大小
func (s *RPCServerPool) ResetAllPoolSize(minSize int, maxSize int) {
	s.MinSize = minSize
	s.MaxSize = maxSize
	s.pool.ResetAllPoolSize(minSize, maxSize)
}

//Close 关闭连接池
func (s *RPCServerPool) Close() {
	s.pool.Close()
}

//Register 注册服务列表
func (s *RPCServerPool) Register(svs map[string]string) {
	//标记不能使用的服务
	servers := s.servers.GetAll()
	for ip := range servers {
		if _, ok := svs[ip]; !ok {
			s.servers.Delete(ip)
			go func(ip string) {
				defer s.recover()
				s.pool.UnRegister(ip)
				removeWorker(ip)
			}(ip)
		}
	}
	//*
	//添加可以使用使用的服务
	for ip := range svs {
		if _, ok := servers[ip]; !ok {
			go func(ip string) {
				defer s.recover()
				err := s.pool.Register(ip, newRPCClientFactory(ip, s.loggerName), s.MinSize, s.MaxSize)
				if err != nil {
					s.Log.Error(err)
					return
				}
				s.servers.Set(ip, &rpcServerService{IP:ip, Status: true}) //set ->add
			}(ip)
		}
	}
}

func (p *RPCServerPool) Request(group string, svName string, input string, session string) (result string, err error) {
	defer p.recover()
	if strings.EqualFold(group, "") {
		err = errors.New("not find rpc server and name cant be nil" + p.loggerName + "@" + p.domain + ".rpc.pool")
		return
	}
	execute := 0
START:
	execute++
	if execute > p.MaxRetry {
		err = fmt.Errorf("cant connect to rpc server(%s@%s.rpc.pool):%s/%s,%s", p.loggerName, p.domain, group, svName, err)
		return
	}
	o, err := p.pool.Get(group)
	if err != nil {
		err = fmt.Errorf("not find rpc server(%s@%s.rpc.pool):%s/%s,%s", p.loggerName, p.domain, group, svName, err)
		return
	}
	obj := o.(*RPCClient)
	err = obj.Open()
	defer obj.Close()
	if err != nil {
		p.Log.Error("当前服务不可用:", p.loggerName, svName, err)
		p.pool.Unusable(svName, obj)
		goto START
	}
	defer p.pool.Recycle(group, o)
	result, err = obj.Request(svName, input, session)
	return
}
func (p *RPCServerPool) Send(group string, svName string, input string, data []byte) (result string, err error) {
	defer p.recover()
	if strings.EqualFold(group, "") {
		err = errors.New("not find rpc server and name cant be nil")
		return
	}

	o, err := p.pool.Get(group)
	if err != nil {
		err = errors.New("not find rpc server")
		return
	}
	obj := o.(*RPCClient)
	result, err = obj.Send(svName, input, data)
	if err != nil {
		p.pool.Unusable(svName, obj)
	} else {
		p.pool.Recycle(group, o)
	}
	return
}

//GetSnap 获取当前RPC客户端的连接池快照信息
func (p *RPCServerPool) GetSnap() pool.ObjectPoolSnap {
	return p.pool.GetSnap()
}

func (p *RPCServerPool) Get(group string, svName string, input string) (result []byte, err error) {
	defer p.recover()
	if strings.EqualFold(group, "") {
		err = errors.New("not find rpc server and name cant be nil")
		return
	}

	o, err := p.pool.Get(group)
	if err != nil {
		err = fmt.Errorf("not find rpc server:%s", err)
		return
	}
	obj := o.(*RPCClient)
	result, err = obj.Get(svName, input)
	if err != nil {
		p.pool.Unusable(svName, obj)
	} else {
		p.pool.Recycle(group, o)
	}
	return
}
func (n *RPCServerPool) recover() {
	if r := recover(); r != nil {
		n.Log.Fatal(r, string(debug.Stack()))
	}
}
