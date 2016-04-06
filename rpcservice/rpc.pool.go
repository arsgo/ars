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

func (p *RPCServerPool) executeRequest(svname, cmd, input string) (result string, err error) {
	defer func() {
		if ex := recover(); ex != nil {
			err = ex.(error)
		}
	}()

	o, ex := p.pool.Get(svname)
	if ex != nil {
		o.Fatal()
		return "", errors.New("not find rc server")
	}
	if !o.Check() {
		return "", errors.New("not find available rc server")
	}
	defer p.pool.Recycle(svname, o)
	obj := o.(*rpcClient)
	return obj.Request(cmd, input)
}
func (p *RPCServerPool) executeSend(svname, cmd, input string, data []byte) (result string, err error) {
	defer func() {
		if ex := recover(); ex != nil {
			err = ex.(error)
		}
	}()

	o, ex := p.pool.Get(svname)
	if ex != nil {
		o.Fatal()
		return "", errors.New("not find rc server")
	}
	if !o.Check() {
		return "", errors.New("not find available rc server")
	}
	defer p.pool.Recycle(svname, o)
	obj := o.(*rpcClient)
	return obj.Send(cmd, input, data)
}

//Request 执行Request请求
func (p *RPCServerPool) Request(name string, input string) (result string, err error) {
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
		result, err = p.executeRequest(sv, name, input)
		if err == nil {
			p.Log.Infof("->%d:break", index)
			break
		}
	}
	return

}

//Send 发送Send请求
func (p *RPCServerPool) Send(name string, input string, data []byte) (result string, err error) {
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
		result, err = p.executeSend(sv, name, input, data)
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
