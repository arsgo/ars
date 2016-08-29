package script

import lua "github.com/yuin/gopher-lua"

func (s *ScriptPool) moduleCreateMem(ls *lua.LState) int {
	name := ls.CheckString(1)
	client, err := s.NewMemcached(name)
	return pushValues(ls, client, err)
}
