package cluster

import (
	"log"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/script"
)

func (a *spServer) bindGlobalLibs() (funs map[string]interface{}) {
	funs = map[string]interface{}{
		"print":         a.Log.Info,
		"printf":        a.Log.Infof,
		"error":         a.Log.Error,
		"errorf":        a.Log.Errorf,
		"NewMQProducer": a.NewMQProducer,
		"NewElastic":    a.NewElastic,
	}
	return
}

func NewScript(p *spServer) *spScriptEngine {
	var err error
	pool := script.NewLuaPool()
	pool.RegisterLibs(p.bindGlobalLibs())
	en := &spScriptEngine{script: pool, provider: p}
	en.Log, err = logger.New("app script", true)
	if err != nil {
		log.Println(err)
	}
	return en
}
