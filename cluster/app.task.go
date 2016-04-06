package cluster

import (
	"strings"

	"github.com/colinyl/ars/scheduler"
)

func (a *appServer) BindTask(config *AppConfig, err error) error {
	scheduler.Stop()
	if len(config.Auto) == 0 {
		return nil
	}
	for _, v := range config.Auto {		
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
	scheduler.Start()
	return nil
}
