package rpcservice

import (
	"errors"
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
			go s.pool.Register(ip, newRPCClientFactory(ip), 1)
			s.servers[ip] = &rpcServerService{IP: ip, Status: true}
		}
	}
}

func (p *RPCServerPool) Request(group string, svName string, input string) (result string, err error) {
	defer func() {
		if ex := recover(); ex != nil {
			err = ex.(error)
		}
	}()

	o, ex := p.pool.Get(group)
	if ex != nil {
		o.Fatal()
		return "", errors.New("not find rc server")
	}
	if !o.Check() {
		return "", errors.New("not find available rc server")
	}
	defer p.pool.Recycle(group, o)
	obj := o.(*rpcClient)
	return obj.Request(svName, input)
}
func (p *RPCServerPool) Send(group string, svName string, input string, data []byte) (result string, err error) {
	defer func() {
		if ex := recover(); ex != nil {
			err = ex.(error)
		}
	}()

	o, ex := p.pool.Get(group)
	if ex != nil {
		o.Fatal()
		return "", errors.New("not find rc server")
	}
	if !o.Check() {
		return "", errors.New("not find available rc server")
	}
	defer p.pool.Recycle(group, o)
	obj := o.(*rpcClient)
	return obj.Send(svName, input, data)
}

func (p *RPCServerPool) Get(group string, svName string, input string) (result []byte, err error) {
	defer func() {
		if ex := recover(); ex != nil {
			err = ex.(error)
		}
	}()

	o, ex := p.pool.Get(group)
	if ex != nil {
		o.Fatal()
		return make([]byte,0), errors.New("not find rc server")
	}
	if !o.Check() {
		return make([]byte,0), errors.New("not find available rc server")
	}
	defer p.pool.Recycle(group, o)
	obj := o.(*rpcClient)
	return obj.Get(svName, input)
}
//Request 执行Request请求
func (p *RPCServerPool) Request1(name string, input string) (result string, err error) {
	p.lk.Lock()
	defer p.lk.Unlock()
	if len(p.servers) == 0 {
		return "", errors.New("not find rc server")
	}
	p.Log.Infof("servers:%d", len(p.servers))
	var index int
	for sv, server := range p.servers {
		index++
		p.Log.Infof("->%d", index)
		if !server.Status {
			err = errors.New("not find available rc server")
			continue
		}
		result, err = p.Request(sv, name, input)
		if err == nil {
			p.Log.Infof("->%d:break", index)
			break
		}
	}
	return

}

//Send 发送Send请求
func (p *RPCServerPool) Send1(name string, input string, data []byte) (result string, err error) {
	p.lk.Lock()
	defer p.lk.Unlock()
	if len(p.servers) == 0 {
		return "", errors.New("not find rc server")
	}
	p.Log.Infof("servers:%d", len(p.servers))
	var index int
	for sv, server := range p.servers {
		index++
		p.Log.Infof("->%d", index)
		if !server.Status {
			err = errors.New("not find available rc server")
			continue
		}
		result, err = p.Send(sv, name, input, data)
		if err == nil {
			p.Log.Infof("->%d:break", index)
			break
		}
	}
	return
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
