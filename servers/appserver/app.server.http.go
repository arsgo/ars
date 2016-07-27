package main

import (
	"fmt"
	"strings"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/server"
)

func (a *AppServer) BindHttpServer(config *cluster.ServerConfig) {
	if config == nil {
		if a.httpServer != nil {
			a.httpServer.Stop()
		}
		return
	}
	if a.httpServer != nil && !strings.EqualFold(a.httpServer.Address, config.Address) {
		a.httpServer.Stop()
	}

	for _, v := range config.Routes {
		er := a.scriptPool.PreLoad(v.Script, v.MinSize, v.MaxSize)
		if er != nil {
			a.Log.Error("load script error in:", v.Script, ",", er)
		} else {
			a.Log.Info("::load script ", v.Script)
		}
	}
	if config != nil && len(config.Routes) > 0 &&
		strings.EqualFold(config.ServerType, "http") {
		var err error
		a.httpServer, err = server.NewHTTPScriptServer(config.Address, config.Routes, a.scriptPool.Call, a.loggerName)
		if err == nil {
			a.httpServer.Start()
			a.snap.Server = fmt.Sprint(a.snap.ip, a.httpServer.Address)
		} else {
			a.Log.Error(err)
		}
	}
}
