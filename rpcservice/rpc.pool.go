package rpcservice

import (
	"errors"
	"strings"
	"time"
)

//Register 注册服务列表
func (s *RPCServerPool) Register(svs map[string]string) {
	//标记不能使用的服务
	servers := s.servers.GetAll()
	for _, sv := range servers {
		server := sv.(*rpcServerService)
		if _, ok := svs[server.IP]; !ok {
			s.servers.Delete(server.IP)
			go s.pool.UnRegister(server.IP)
		}
	}
	//*
	//添加可以使用使用的服务
	for _, ip := range svs {
		if _, ok := servers[ip]; !ok {
			go func() {
				s.servers.Set(ip, &rpcServerService{IP: ip, Status: true})
				err := s.pool.Register(ip, newRPCClientFactory(ip, s.Log), s.MinSize, s.MaxSize)
				if err != nil {
					s.Log.Error(err)
				}
			}()
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
		err = errors.New("not find rpc server")
		return
	}
	defer p.pool.Recycle(group, o)
	obj := o.(*RPCClient)
	result, err = obj.Request(svName, input)
	if err != nil {
		p.pool.Unusable(svName, obj)
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
	defer p.pool.Recycle(group, o)
	obj := o.(*RPCClient)
	result, err = obj.Send(svName, input, data)
	if err != nil {
		p.pool.Unusable(svName, obj)
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
	defer p.pool.Recycle(group, o)
	obj := o.(*RPCClient)
	result, err = obj.Get(svName, input)
	if err != nil {
		p.pool.Unusable(svName, obj)
	}
	return
}
func (p *RPCServerPool) clearUp() {
	/*p.lk.Lock()
	for k, server := range p.servers {
		if !server.Status && p.pool.Close(server.IP) {
			delete(p.servers, k)
		}
	}
	p.lk.Unlock()*/
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
