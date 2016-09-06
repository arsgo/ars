package script

import (
	"github.com/arsgo/lib4go/security/aes"
	"github.com/arsgo/lib4go/security/base64"
	"github.com/arsgo/lib4go/security/des"
	"github.com/arsgo/lib4go/security/md5"
	"github.com/arsgo/lib4go/security/rsa"
	"github.com/arsgo/lib4go/security/sha1"
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
func (s *ScriptPool) moduleBase64EncodeBytes(ls *lua.LState) int {
	input := ls.CheckUserData(1)
	data := input.Value.([]byte)
	return pushValues(ls, base64.EncodeBytes(data))
}

func (s *ScriptPool) moduleBase64Decode(ls *lua.LState) int {
	input := ls.CheckString(1)
	r, e := base64.Decode(input)
	return pushValues(ls, r, e)
}
func (s *ScriptPool) moduleBase64DecodeBytes(ls *lua.LState) int {
	input := ls.CheckString(1)
	r, e := base64.DecodeBytes(input)
	return pushValues(ls, r, e)
}
func (s *ScriptPool) moduleSha1Encrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	return pushValues(ls, sha1.Encrypt(input))
}
func (s *ScriptPool) moduleRsaEncrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	publicKey := ls.CheckString(2)
	data, err := rsa.Encrypt(input, publicKey)
	return pushValues(ls, data, err)
}
func (s *ScriptPool) moduleRsaDecrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	privateKey := ls.CheckString(2)
	data, err := rsa.Decrypt(input, privateKey)
	return pushValues(ls, data, err)
}
func (s *ScriptPool) moduleRsaMakeSign(ls *lua.LState) int {
	input := ls.CheckString(1)
	privateKey := ls.CheckString(2)
	data, err := rsa.Sign(input, privateKey)
	return pushValues(ls, data, err)
}
func (s *ScriptPool) moduleRsaVerify(ls *lua.LState) int {
	src := ls.CheckString(1)
	sign := ls.CheckString(2)
	pubkey := ls.CheckString(3)
	data, err := rsa.Verify(src, sign, pubkey)
	return pushValues(ls, data, err)
}
