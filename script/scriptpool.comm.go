package script

import (
	"github.com/colinyl/ars/servers/config"
	"github.com/colinyl/lib4go/utility"
	"github.com/yuin/gopher-lua"
)

func (s *ScriptPool) moduleGetGUID(ls *lua.LState) int {
	return pushValues(ls, utility.GetGUID())
}
func (s *ScriptPool) moduleGetLocalIP(ls *lua.LState) int {
	ip, err := config.Get()
	if err != nil {
		return pushValues(ls, "", err)
	}
	return pushValues(ls, ip.IP)

}
