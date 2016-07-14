package rpcproxy

import (
	"github.com/colinyl/lib4go/security/aes"
	"github.com/colinyl/lib4go/security/base64"
	"github.com/colinyl/lib4go/security/des"
	"github.com/colinyl/lib4go/security/md5"
	"github.com/colinyl/lib4go/security/sha1"
	"github.com/yuin/gopher-lua"
)

func (s *ScriptPool) moduleMd5Encrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	return pushValues(ls, md5.Encrypt(input))
}
func (s *ScriptPool) moduleDESEncrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	key := ls.CheckString(2)
	r, e := des.Encrypt(input, key)
	return pushValues(ls, r, e)
}
func (s *ScriptPool) moduleDESDecrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	key := ls.CheckString(2)
	r, e := des.Decrypt(input, key)
	return pushValues(ls, r, e)
}
func (s *ScriptPool) moduleAESEncrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	key := ls.CheckString(2)
	r, e := aes.Encrypt(input, key)
	return pushValues(ls, r, e)
}
func (s *ScriptPool) moduleAESDecrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	key := ls.CheckString(2)
	r, e := aes.Decrypt(input, key)
	return pushValues(ls, r, e)
}
func (s *ScriptPool) moduleBase64Encode(ls *lua.LState) int {
	input := ls.CheckString(1)
	return pushValues(ls, base64.Encode(input))
}
func (s *ScriptPool) moduleBase64Decode(ls *lua.LState) int {
	input := ls.CheckString(1)
	r, e := base64.Decode(input)
	return pushValues(ls, r, e)
}
func (s *ScriptPool) moduleSha1Encrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	return pushValues(ls, sha1.Encrypt(input))
}