package main

import (
	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/scheduler"
)

//BindLocalJobs 绑定本地JOB
func (a *AppServer) BindLocalJobs(jobs []cluster.JobItem) {
	scheduler.Stop()
	if len(jobs) == 0 {
		return
	}
	aliveJob := 0
	for _, v := range jobs {
		if !v.Enable {
			continue
		}
		aliveJob++
		a.Log.Infof("::start local job:[%s] %s", v.Name, v.Script)
		er := a.scriptPool.Pool.PreLoad(v.Script, v.MinSize, v.MaxSize)
		if er != nil {
			a.Log.Error("load task`s script error in:", v.Script, ",", er)
			continue
		}
		scheduler.AddTask(v.Trigger, scheduler.NewTask(v, func(job interface{}) {			
			item := job.(cluster.JobItem)
			a.Log.Infof(" -> run job [%s] %s", item.Name, item.Script)
			_, err := a.scriptPool.Call(item.Script, "{}", item.Params)
			if err != nil {
				a.Log.Error(err)
			}

		}))
	}
	if aliveJob > 0 {
		scheduler.Start()
	}

}
