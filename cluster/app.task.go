package cluster

import (
	"strings"

	"github.com/colinyl/ars/scheduler"
)

func (a *appServer) BindTask(config *AppConfig, err error) error {
	a.resetSnap(a.appServerAddress, config)
	scheduler.Stop()
	for _, v := range config.Tasks {
		scheduler.AddTask(v.Trigger, scheduler.NewTask(v.Script, func(name interface{}) {
			a.Log.Infof("start:%s", name)
			rtvalues, err := a.scriptEngine.pool.Call(name.(string))
			if err != nil {
				a.Log.Error(err)
			} else {
				a.Log.Infof("result:%d,%s", len(rtvalues), strings.Join(rtvalues, ","))
			}
		}))
	}
	if len(config.Jobs) > 0 {
		a.StartJobConsumer(config.Jobs)
	} else {
		a.StopJobServer()
	}
	if len(config.Tasks) > 0 {
		scheduler.Start()
	} else {
		scheduler.Stop()
	}
	if config.Server != nil && len(config.Server.Routes) > 0 && strings.EqualFold(config.Server.ServerType, "http") {
		a.appRoutes = config.Server.Routes
		a.StopHttpAPIServer()
		a.StartHttpAPIServer()
	} else {
		a.StopHttpAPIServer()
	}
	a.monitor.Bind(config.Monitor)

	return nil
}
