package script

import (
	"github.com/arsgo/lib4go/webserver"
	lua "github.com/yuin/gopher-lua"
)

func (s *ScriptPool) moduleContextGetHTTPContext(ls *lua.LState) *webserver.Context {
	context := ls.GetGlobal("__http_context__")
	if context == nil {
		return nil
	}
	data := context.(*lua.LUserData)
	hr := data.Value.(*webserver.Context)
	return hr
}
func (s *ScriptPool) moduleContextGetCookie(ls *lua.LState) int {
	key := ls.CheckString(1)
	context := s.moduleContextGetHTTPContext(ls)
	if context == nil {
		return pushValues(ls)
	}
	ck, err := context.Request.Cookie(key)
	if err != nil {
		return pushValues(ls)
	}
	return pushValues(ls, ck.Value)

}

func (s *ScriptPool) moduleContextSetCookie(ls *lua.LState) int {
	key := ls.CheckString(1)
	value := ls.CheckString(2)
	cookie := ls.GetGlobal("__set_cookie__")
	if cookie == lua.LNil {
		ls.SetGlobal("__set_cookie__", lua.LString(key+"="+value))

	} else {
		ls.SetGlobal("__set_cookie__", lua.LString(cookie.String()+";"+key+"="+value))
	}
	return pushValues(ls)
}
func (s *ScriptPool) moduleContextSetContentType(ls *lua.LState) int {
	value := ls.CheckString(1)
	context := s.moduleContextGetHTTPContext(ls)
	if context == nil {
		return pushValues(ls)
	}
	context.Writer.Header().Set("Content-Type", value)
	return pushValues(ls)
}
func (s *ScriptPool) moduleContexSetCharset(ls *lua.LState) int {
	value := ls.CheckString(1)
	context := s.moduleContextGetHTTPContext(ls)
	if context == nil {
		return pushValues(ls)
	}
	context.Writer.Header().Set("Charset", value)
	return pushValues(ls)
}
func (s *ScriptPool) moduleContexSetHeader(ls *lua.LState) int {
	key := ls.CheckString(1)
	value := ls.CheckString(2)
	context := s.moduleContextGetHTTPContext(ls)
	if context == nil {
		return pushValues(ls)
	}
	context.Writer.Header().Set(key, value)
	return pushValues(ls)
}

func (s *ScriptPool) moduleContexSetRaw(ls *lua.LState) int {
	ls.SetGlobal("__raw__", lua.LString("true"))
	return pushValues(ls)
}
