package rpcservice

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
	"time"

	"github.com/arsgo/ars/servers/config"
	"github.com/arsgo/lib4go/concurrent"
	"github.com/arsgo/lib4go/logger"
	"github.com/arsgo/lib4go/pool"
)

type rpcServerService struct {
	Status bool
	IP     string
}

type RPCServerPool struct {
	pool       *pool.ObjectPool
	servers    *concurrent.ConcurrentMap
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
	pl := &RPCServerPool{MinSize: minSize, MaxSize: maxSize, MaxRetry: 3, loggerName: loggerName, domain: conf.Domain}
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
				s.servers.Set(ip, &rpcServerService{IP: ip, Status: true}) //set ->add
			}(ip)
		}
	}
}

//Request 发送request请求
func (p *RPCServerPool) Request(group string, svName string, input string, session string, timeout time.Duration) (result string, err error) {
	defer p.recover()
	//defer base.RunTime("rpc request total", time.Now())
	if strings.EqualFold(group, "") {
		err = errors.New("not find rpc server and name cant be nil" + p.loggerName + "@" + p.domain + ".rpc.pool")
		return
	}
	execute := 0
START:
	execute++
	if execute > p.MaxRetry {
		err = fmt.Errorf("cant connect to rpc server(%s@%s.rpc.pool):%s/%s,%v", p.loggerName, p.domain, group, svName, err)
		return
	}
	o, err := p.pool.Get(group)
	if err != nil {
		err = fmt.Errorf("not find rpc server(%s@%s.rpc.pool):%s/%s,[%v]", p.loggerName, p.domain, group, svName, err)
		return
	}
	obj := o.(*RPCClient)
	err = obj.OpenTimeout(time.Second)
	if err != nil {
		p.Log.Error("当前服务不可用:", p.loggerName, svName, err)
		p.pool.Unusable(group, obj)
		goto START
	}
	obj.SetRWTimeout(timeout)
	result, err = obj.Request(svName, input, session, timeout)
	p.pool.Recycle(group, o)
	obj.Close()
	return
}

//Request 发送request请求
func (p *RPCServerPool) Request2(group string, svName string, input string, session string, timeout time.Duration) (result string, err error) {
	defer p.recover()
	//defer base.RunTime("rpc request total", time.Now())
	if strings.EqualFold(group, "") {
		err = errors.New("not find rpc server and name cant be nil" + p.loggerName + "@" + p.domain + ".rpc.pool")
		return
	}
	execute := 0
START:
	execute++
	if execute > p.MaxRetry {
		err = fmt.Errorf("cant connect to rpc server(%s@%s.rpc.pool):%s/%s,%v", p.loggerName, p.domain, group, svName, err)
		return
	}
	o, err := p.pool.Get(group)
	if err != nil {
		err = fmt.Errorf("not find rpc server(%s@%s.rpc.pool):%s/%s,[%v]", p.loggerName, p.domain, group, svName, err)
		return
	}
	obj := o.(*RPCClient)
	//open,close
	if !obj.Available() {
		p.Log.Error("当前服务不可用:", p.loggerName, svName, err)
		//p.pool.Unusable(group, obj)
		p.pool.Recycle(group, o)
		goto START
	}
	obj.SetRWTimeout(timeout)
	result, err = obj.Request(svName, input, session, timeout)
	p.pool.Recycle(group, o)
	return
}
func (p *RPCServerPool) Send(group string, svName string, input string, data []byte) (result string, err error) {
	defer p.recover()
	return
}

//GetSnap 获取当前RPC客户端的连接池快照信息
func (p *RPCServerPool) GetSnap() pool.ObjectPoolSnap {
	return p.pool.GetSnap()
}

func (p *RPCServerPool) Get(group string, svName string, input string) (result []byte, err error) {
	defer p.recover()
	return
}
func (n *RPCServerPool) recover() {
	if r := recover(); r != nil {
		n.Log.Fatal(r, string(debug.Stack()))
	}
}
