package rpcproxy

import "github.com/yuin/gopher-lua"

func (s *ScriptPool) bindModules() (r map[string]map[string]lua.LGFunction) {
	r = map[string]map[string]lua.LGFunction{
		"common": map[string]lua.LGFunction{
			"getGuid": s.moduleGetGUID,
			"getIP":   s.moduleGetLocalIP,
		},
		"rpc": map[string]lua.LGFunction{
			"request": s.moduleRPCRequest,
		},
		"md5": map[string]lua.LGFunction{
			"encrypt": s.moduleMd5Encrypt,
		},
		"des": map[string]lua.LGFunction{
			"encrypt": s.moduleDESEncrypt,
			"decrypt": s.moduleDESDecrypt,
		},
		"aes": map[string]lua.LGFunction{
			"encrypt": s.moduleAESEncrypt,
			"decrypt": s.moduleAESDecrypt,
		},
		"base64": map[string]lua.LGFunction{
			"encode": s.moduleBase64Encode,
			"decode": s.moduleBase64Decode,
		},
		"sha1": map[string]lua.LGFunction{
			"encrypt": s.moduleSha1Encrypt,
		},
	}
	return
}
