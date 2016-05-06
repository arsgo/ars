package rpcservice

import (
	"errors"
	"strings"
	"time"
)

//Register 注册服务列表
func (s *RPCServerPool) Register(svs map[string]string) {
	s.lk.Lock()
	defer s.lk.Unlock()
	//标记不能使用的服务
	for _, server := range s.servers {
		if _, ok := svs[server.IP]; !ok {
			server.Status = false
		}
	}
	//添加可以使用使用的服务
	for _, ip := range svs {
		if sv, ok := s.servers[ip]; !ok || !sv.Status {
			s.pool.UnRegister(ip)
			go func() {
				err := s.pool.Register(ip, newRPCClientFactory(ip, s.Log), 10)
				if err != nil {
					s.Log.Error(err)
				}
			}()
			s.servers[ip] = &rpcServerService{IP: ip, Status: true}
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
		o.Fatal()
		err = errors.New("not find rpc server")
		return
	}
	defer p.pool.Recycle(group, o)
	if !o.Check() {
		err = errors.New("not find available rpc server")
		return
	}
	obj := o.(*RPCClient)
	result, err = obj.Request(svName, input)
	if err != nil {
		p.Log.Error(err)
		obj.Fatal()
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
		o.Fatal()
		err = errors.New("not find rpc server")
		return
	}
	defer p.pool.Recycle(group, o)
	if !o.Check() {
		err = errors.New("not find available rpc server")
		return
	}
	obj := o.(*RPCClient)
	return obj.Send(svName, input, data)
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
		o.Fatal()
		err = errors.New("not find rpc server")
		return
	}
	defer p.pool.Recycle(group, o)
	if !o.Check() {
		err = errors.New("not find available rpc server")
		return
	}
	obj := o.(*RPCClient)
	return obj.Get(svName, input)
}
func (p *RPCServerPool) clearUp() {
	p.lk.Lock()
	for k, server := range p.servers {
		if !server.Status && p.pool.Close(server.IP) {
			delete(p.servers, k)
		}
	}
	p.lk.Unlock()
}

func (p *RPCServerPool) autoClearUp() {
	timepk := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-timepk.C:
			{
				p.clearUp()
			}
		}
	}
}
