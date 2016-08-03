package script

import "github.com/yuin/gopher-lua"

func (s *ScriptPool) bindGlobal() (r map[string]lua.LGFunction) {
	r = map[string]lua.LGFunction{
		"print":  s.globalInfo,
		"printf": s.globalInfof,
		"error":  s.globalError,
		"errorf": s.globalErrorf,
		"sleep":  s.globalSleep,
	}
	return
}
