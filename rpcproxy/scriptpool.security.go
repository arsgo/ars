package rpcproxy

import (
	"github.com/colinyl/lib4go/security/aes"
	"github.com/colinyl/lib4go/security/base64"
	"github.com/colinyl/lib4go/security/des"
	"github.com/colinyl/lib4go/security/md5"
	"github.com/colinyl/lib4go/security/sha1"
	"github.com/colinyl/lib4go/security/weixin"
)

//BindSecurity 安全绑定函数
type BindSecurity struct {
}

//NewBindSecurity 创建用于加解密的绑定函数
func (s *ScriptPool) NewBindSecurity() BindSecurity {
	return BindSecurity{}
}

//MD5 MD5加密
func (b BindSecurity) MD5(input string) string {
	return md5.Encrypt(input)
}

//DESEncrypt DES加密
func (b BindSecurity) DESEncrypt(input string, key string) (string, error) {
	return des.Encrypt(input, key)
}

//DESDecrypt DES解密
func (b BindSecurity) DESDecrypt(input string, key string) (string, error) {
	return des.Decrypt(input, key)
}

//AESEncrypt AES加密
func (b BindSecurity) AESEncrypt(input string, key string) (string, error) {
	return aes.Encrypt(input, key)
}

//AESDecrypt AES解密
func (b BindSecurity) AESDecrypt(input string, key string) (string, error) {
	return aes.Decrypt(input, key)
}

//Base64Encode Base64编码
func (b BindSecurity) Base64Encode(input string) string {
	return base64.Encode(input)
}

//Base64Decode Base64解码
func (b BindSecurity) Base64Decode(input string) (string, error) {
	return base64.Decode(input)
}

//SHA1Encrypt SHA1加密
func (b BindSecurity) SHA1Encrypt(input string) string {
	return sha1.Encrypt(input)
}

//MD5Encrypt md5加密
func (b BindSecurity) MD5Encrypt(input string) string {
	return md5.Encrypt(input)
}

//WXDecrypt 微信报文解密
func (b BindSecurity) WXDecrypt(content string) (r string, err error) {
	return weixin.Decrypt(content)
}

//WXEncrypt 微信报文加密
func (b BindSecurity) WXEncrypt(fromUserName, toUserName, content, nonce, timestamp string) (r string, err error) {
	return weixin.Encrypt(fromUserName, toUserName, content, nonce, timestamp)
}

//WXMakeSign 微信生成签名
func (b BindSecurity) WXMakeSign(timestamp, nonce, msgEncrypt string) string {
	return weixin.MakeSign(timestamp, nonce, msgEncrypt)
}
