package rpcservice

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/colinyl/lib4go/concurrent"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/pool"
)

type rpcServerService struct {
	Status bool
	IP     string
}

type RPCServerPool struct {
	pool    *pool.ObjectPool
	servers concurrent.ConcurrentMap
	Log     *logger.Logger
	MinSize int
	MaxSize int
}

func NewRPCServerPool() *RPCServerPool {
	var err error
	pl := &RPCServerPool{}
	pl.pool = pool.New()
	pl.servers = concurrent.NewConcurrentMap()
	pl.Log, err = logger.New("rc server", true)
	if err != nil {
		log.Println(err)
	}
	return pl
}

//Register 注册服务列表
func (s *RPCServerPool) Register(svs map[string]string) {
	//标记不能使用的服务
	servers := s.servers.GetAll()
	for ip := range servers {
		if _, ok := svs[ip]; !ok {
			s.servers.Delete(ip)
			go s.pool.UnRegister(ip)
		}
	}
	//*
	//添加可以使用使用的服务
	for ip := range svs {
		if _, ok := servers[ip]; !ok {
			go func(ip string) {
				err := s.pool.Register(ip, newRPCClientFactory(ip, s.Log), s.MinSize, s.MaxSize)
				if err != nil {
					s.Log.Error(err)
				}
				s.servers.Set(ip, &rpcServerService{IP: ip, Status: true})

			}(ip)
		}
	}
}

func (p *RPCServerPool) Request(group string, svName string, input string) (result string, err error) {
	defer func() {
		if ex := recover(); ex != nil {
		}
	}()
	if strings.EqualFold(group, "") {
		err = errors.New("not find rpc server and name cant be nil")
		return
	}
	o, err := p.pool.Get(group)
	if err != nil {
		err = fmt.Errorf("not find rpc server:%s/%s", group, svName)
		return
	}
	obj := o.(*RPCClient)
	result, err = obj.Request(svName, input)
	if err != nil {
		p.pool.Unusable(svName, obj)
	} else {
		p.pool.Recycle(group, o)
	}
	return
}
func (p *RPCServerPool) Send(group string, svName string, input string, data []byte) (result string, err error) {
	defer func() {
		if ex := recover(); ex != nil {
			//err = ex.(error)
		}
	}()
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

func (p *RPCServerPool) Get(group string, svName string, input string) (result []byte, err error) {
	defer func() {
		if ex := recover(); ex != nil {
			//err = ex.(error)
		}
	}()
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
	result, err = obj.Get(svName, input)
	if err != nil {
		p.pool.Unusable(svName, obj)
	} else {
		p.pool.Recycle(group, o)
	}
	return
}
