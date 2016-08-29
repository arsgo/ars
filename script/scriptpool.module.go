package script

import "github.com/yuin/gopher-lua"

func (s *ScriptPool) bindModules() (r map[string]map[string]lua.LGFunction) {
	r = map[string]map[string]lua.LGFunction{
		"mq": map[string]lua.LGFunction{
			"send": s.moduleMQProducerSend,
		},
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
