package main

import (
	"fmt"
	"strings"

	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/ars/server"
)

func (a *AppServer) BindAPIServer(config *cluster.ServerConfig) {
	if config == nil {
		a.closeAPIServer("api server 配置文件为空")
		return
	}
	if !strings.EqualFold(config.ServerType, "http") {
		a.closeAPIServer(fmt.Sprintf("服务器类型:%s 不支持", config.ServerType))
		return
	}

	if len(config.Routes) == 0 {
		a.closeAPIServer("api server 路由未配置")
		return
	}

	if config.Disable {
		a.closeAPIServer("api server 配置未启用")
		return
	}

	needCreate := a.apiServer == nil
	needStart := a.apiServer != nil && !strings.EqualFold(strings.Trim(a.apiServer.Address, ":"), config.Address)
	if needStart {
		a.Log.Info("api server 端口号已变更，停止api server:", a.apiServer.Address, config.Address)
		a.apiServer.Stop()
	}

	for _, v := range config.Routes {
		er := a.scriptPool.PreLoad(v.Script, v.MinSize, v.MaxSize)
		if er != nil {
			a.Log.Error("脚本加载失败:", v.Script, ",", er)
		}
	}

	if needCreate {
		var err error
		a.apiServer, err = server.NewHTTPScriptServer(config.Address, config.Routes, a.scriptPool.Call, a.loggerName, a.apiServerCollector)
		if err != nil {
			a.Log.Error("api server 启动失败:", err)
			return
		}
		needStart = true
	}
	if needStart || !a.apiServer.Available {
		if err := a.apiServer.Start(); err != nil {
			a.clusterClient.CloseNode(a.snap.path)
			return
		}
	}
	a.snap.port = a.apiServer.Address
	a.snap.Server = fmt.Sprint(a.conf.IP, a.apiServer.Address)
	a.createAPIServer()
}

func (a *AppServer) createAPIServer() {
	if a.snap.path == "" {
		a.snap.path, _ = a.clusterClient.CreateAppServer(a.snap.port, a.snap.GetSnap())
	}
}

func (a *AppServer) closeAPIServer(msg string) bool {

	if a.apiServer != nil && a.apiServer.Available {
		a.Log.Info(" -> ", msg, ", 停止api server")
		if !strings.EqualFold(a.snap.path, "") {
			a.clusterClient.CloseNode(a.snap.path)
		}
		a.apiServer.Stop()
	} else {
		a.Log.Info(" -> ", msg)
	}
	a.snap.path = ""
	return a.apiServer != nil
}
