package cluster

import (
	"errors"
	"log"
	"strings"

	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/lua"
)

type spScriptEngine struct {
	script   *lua.LuaPool
	provider *spServer
	Log      *logger.Logger
}

func (s *spScriptEngine) Request(cmd string, input string) (string, error) {
	s.Log.Infof("recv request:%s", cmd)
	s.provider.lk.Lock()
	svs, ok := s.provider.services.services[cmd]
	s.provider.lk.Unlock()
	if !ok {
		return "", errors.New("not suport")
	}
	path := svs.Script
	values, err := s.script.Call(path, input)
	return strings.Join(values, ","), err
}
func (s *spScriptEngine) Send(cmd string, input string, data []byte) (string, error) {
	s.provider.lk.Lock()
	svs, ok := s.provider.services.services[cmd]
	s.provider.lk.Unlock()

	if !ok {
		return "", errors.New("not suport")
	}
	path := svs.Script
	values, err := s.script.Call(path, input)
	return strings.Join(values, ","), err
}
func (s *spScriptEngine) Get(cmd string, input string) ([]byte, error) {
	s.provider.lk.Lock()
	svs, ok := s.provider.services.services[cmd]
	s.provider.lk.Unlock()

	if !ok {
		return make([]byte, 0), errors.New("not suport")
	}
	path := svs.Script
	values, err := s.script.Call(path, input)
	return []byte(strings.Join(values, ",")), err
}

func NewScript(p *spServer) *spScriptEngine {
	var err error
	en := &spScriptEngine{script: lua.NewLuaPool(), provider: p}
	en.Log, err = logger.New("app script", true)
	if err != nil {
		log.Println(err)
	}
	return en
}

func (d *spServer) StartRPC() {
	address := rpcservice.GetLocalRandomAddress()
	d.Port = address
	d.dataMap.Set("port", d.Port)
	rpcServer := rpcservice.NewRPCServer(address, NewScript(d))
	rpcServer.Serve()
}
