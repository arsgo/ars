package script

import (
	"fmt"
	"time"

	"github.com/arsgo/lib4go/logger"
	"github.com/yuin/gopher-lua"
)

func (s *ScriptPool) globalGetParams(ls *lua.LState) (params []interface{}) {
	c := ls.GetTop()
	params = make([]interface{}, 0, c)
	for i := 1; i <= c; i++ {
		t := ls.Get(i).Type().String()
		if t == "userdata" {
			params = append(params, fmt.Sprintf("%+v", ls.CheckUserData(i).Value))
		} else {
			params = append(params, ls.Get(i).String())
		}
	}
	return
}
func (s *ScriptPool) globalGetLogger(ls *lua.LState) (lg logger.ILogger, err error) {
	loggerName := ls.GetGlobal("__logger_name__").String()
	sessionID := ls.GetGlobal("__session__").String()
	lg, err = logger.NewSession(loggerName, sessionID)
	return
}
func (s *ScriptPool) globalInfo(ls *lua.LState) int {
	params := s.globalGetParams(ls)
	if len(params) == 0 {
		return pushValues(ls)
	}
	lg, err := s.globalGetLogger(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	lg.Info(params...)
	return pushValues(ls)
}
func (s *ScriptPool) globalInfof(ls *lua.LState) int {
	params := s.globalGetParams(ls)
	if len(params) <= 1 {
		return pushValues(ls)
	}
	lg, err := s.globalGetLogger(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	lg.Infof(params[0].(string), params[1:]...)
	return pushValues(ls)
}

func (s *ScriptPool) globalError(ls *lua.LState) int {
	params := s.globalGetParams(ls)
	if len(params) == 0 {
		return pushValues(ls)
	}
	lg, err := s.globalGetLogger(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	lg.Error(params...)
	return pushValues(ls)
}
func (s *ScriptPool) globalErrorf(ls *lua.LState) int {
	params := s.globalGetParams(ls)
	if len(params) <= 1 {
		return pushValues(ls)
	}
	lg, err := s.globalGetLogger(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	lg.Errorf(params[0].(string), params[1:]...)
	return pushValues(ls)
}
func (s *ScriptPool) globalSleep(ls *lua.LState) int {
	second := ls.CheckInt(1)
	time.Sleep(time.Second * time.Duration(second))
	return pushValues(ls)
}
