package cluster

import (
	"strings"

	"github.com/colinyl/ars/scheduler"
)

func (a *appServer) BindTask(config *AppConfig, err error) error {
	scheduler.Stop()
	for _, v := range config.Tasks {
		scheduler.AddTask(v.Trigger, scheduler.NewTask(v.Script, func(name string) {
			a.Log.Infof("start:%s", name)
			rtvalues, err := a.scriptEngine.pool.Call(name, v.Input)
			if err != nil {
				a.Log.Error(err)
			} else {
				a.Log.Infof("result:%s", strings.Join(rtvalues, ","))
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
	if len(config.ScriptServer) > 0 {
		a.scriptServer = config.ScriptServer
		a.StopAPIServer()
		a.StartAPIServer()
	} else {
		a.StopAPIServer()
	}

	return nil
}
