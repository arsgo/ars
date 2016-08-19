package main

import (
	"fmt"
	"strings"

	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/ars/server"
)

func (a *AppServer) BindHttpServer(config *cluster.ServerConfig) {
	if config == nil {
		if a.httpServer != nil {
			a.Log.Info("配置文件为空，停止api server")
			a.httpServer.Stop()
		}
		return
	}

	if len(config.Routes) == 0 || !strings.EqualFold(config.ServerType, "http") {
		if a.httpServer != nil {
			a.Log.Info("路由信息未配置或服务器类型错误，停止api server")
			a.httpServer.Stop()
		}
		return
	}

	needCreate := a.httpServer == nil
	needStart := a.httpServer != nil && !strings.EqualFold(strings.Trim(a.httpServer.Address, ":"), config.Address)
	if needStart {
		a.Log.Info("端口号已变更，停止api server:", a.httpServer.Address, config.Address)
		a.httpServer.Stop()
	}

	for _, v := range config.Routes {
		er := a.scriptPool.PreLoad(v.Script, v.MinSize, v.MaxSize)
		if er != nil {
			a.Log.Error("脚本加载失败:", v.Script, ",", er)
		}
	}

	if needCreate {
		var err error
		a.httpServer, err = server.NewHTTPScriptServer(config.Address, config.Routes, a.scriptPool.Call, a.loggerName)
		if err != nil {
			a.Log.Error("http server 启动失败:", err)
			return
		}
		needStart = true
	}
	if needStart {
		a.httpServer.Start()
	}
	a.snap.port = a.httpServer.Address
	a.snap.Server = fmt.Sprint(a.snap.ip, a.httpServer.Address)
}
