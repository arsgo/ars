package cluster

import (
	"errors"
	"fmt"
	"log"
	"strings"

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
		return getErrorResult("500", fmt.Sprintf("not support service %s", cmd)), nil
	}
	path := svs.Script
	values, err := s.script.Call(path, input)
	if err != nil {
		return getErrorResult("500", err.Error()), nil
	}
	s.Log.Info(values[0])
	return getDataResult(strings.Join(values, ",")), nil
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
