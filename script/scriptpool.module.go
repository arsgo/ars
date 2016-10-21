package script

import "github.com/yuin/gopher-lua"

func (s *ScriptPool) bindModules() (r map[string]map[string]lua.LGFunction) {
	r = map[string]map[string]lua.LGFunction{
		"mq": map[string]lua.LGFunction{
			"send": s.moduleMQProducerSend,
		},
		"context": map[string]lua.LGFunction{
			"set_cookie":       s.moduleContextSetCookie,
			"get_cookie":       s.moduleContextGetCookie,
			"set_charset":      s.moduleContexSetCharset,
			"set_header":       s.moduleContexSetHeader,
			"set_content_type": s.moduleContextSetContentType,
			"set_raw":          s.moduleContexSetRaw,
			"get_guid":         s.moduleGetGUID,
			"get_ip":           s.moduleGetLocalIP,
		},
		"common": map[string]lua.LGFunction{
			"get_guid": s.moduleGetGUID,
			"get_ip":   s.moduleGetLocalIP,
		},
		"rpc": map[string]lua.LGFunction{
			"request":       s.moduleRPCRequest,
			"async_request": s.moduleRPCAsyncRequest,
			"get_result":    s.moduleRPCGetResult,
		},
		"url": map[string]lua.LGFunction{
			"encode": s.moduleURLEncode,
			"decode": s.moduleURLDecode,
		},
		"encoding": map[string]lua.LGFunction{
			"convert": s.moduleEncodingConvert,
		},
		"unicode": map[string]lua.LGFunction{
			"encode": s.moduleUnicodeEncode,
			"decode": s.moduleUnicodeDecode,
		},
		"md5": map[string]lua.LGFunction{
			"encrypt": s.moduleMd5Encrypt,
		},
		"des": map[string]lua.LGFunction{
			"encrypt":    s.moduleDESEncrypt,
			"decrypt":    s.moduleDESDecrypt,
			"qx_encrypt": s.moduleQXDESEncrypt,
			"qx_decrypt": s.moduleQXDESDecrypt,
		},
		"aes": map[string]lua.LGFunction{
			"encrypt": s.moduleAESEncrypt,
			"decrypt": s.moduleAESDecrypt,
		},
		"base64": map[string]lua.LGFunction{
			"encode":       s.moduleBase64Encode,
			"decode":       s.moduleBase64Decode,
			"encode_bytes": s.moduleBase64EncodeBytes,
			"decode_bytes": s.moduleBase64DecodeBytes,
		},
		"rsa": map[string]lua.LGFunction{
			"encrypt":   s.moduleRsaEncrypt,
			"decrypt":   s.moduleRsaDecrypt,
			"verify":    s.moduleRsaVerify,
			"make_sign": s.moduleRsaMakeSign,
		},
		"sha1": map[string]lua.LGFunction{
			"encrypt": s.moduleSha1Encrypt,
		},
		"memcached": map[string]lua.LGFunction{
			"new": s.moduleCreateMem,
		},
		"report": map[string]lua.LGFunction{
			"success": s.moduleReportSuccess,
			"error":   s.moduleReportError,
			"failed":  s.moduleReportFaild,
			"juge":    s.moduleReportJuge,
		},
	}
	return
}
