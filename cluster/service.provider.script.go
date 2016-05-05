package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/script"
	"github.com/colinyl/lib4go/utility"
)

type spScriptEngine struct {
	script   *script.LuaPool
	provider *spServer
	Log      *logger.Logger
}
type scriptInputArgs struct {
	Input  json.RawMessage `json:"input"`
	Params json.RawMessage `json:"params"`
}

func getScriptInputArgs(input string, params string) string {
	args := scriptInputArgs{}
	args.Input = []byte(input)
	text, err := utility.GetParams(params)
	if err == nil {
		args.Params = []byte(text)
	}
	buffer, _ := json.Marshal(&args)
	return string(buffer)
}

func (s *spScriptEngine) Request(cmd string, input string) (string, error) {
	s.Log.Infof("recv request:%s", cmd)
	s.provider.lk.Lock()
	svs, ok := s.provider.services.services[cmd]
	s.provider.lk.Unlock()
	if !ok || !s.checkParam(svs, "request") {
		return getErrorResult("500", fmt.Sprintf("not support service %s", cmd)), nil
	}
	values, err := s.script.Call(svs.Script, getScriptInputArgs(input, svs.Params))
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
	if !ok || !s.checkParam(svs, "send") {
		return "", errors.New("not suport")
	}
	values, err := s.script.Call(svs.Script, getScriptInputArgs(input, svs.Params))
	return strings.Join(values, ","), err
}
func (s *spScriptEngine) Get(cmd string, input string) ([]byte, error) {
	s.provider.lk.Lock()
	svs, ok := s.provider.services.services[cmd]
	s.provider.lk.Unlock()

	if !ok || !s.checkParam(svs, "get") {
		return make([]byte, 0), errors.New("not suport")
	}
	values, err := s.script.Call(svs.Script, getScriptInputArgs(input, svs.Params))
	return []byte(strings.Join(values, ",")), err
}

func (s *spScriptEngine) checkParam(v spService, method string) bool {
	return strings.EqualFold(strings.ToLower(v.Type), "rpc") && strings.EqualFold(strings.ToLower(v.Method), method)
}
