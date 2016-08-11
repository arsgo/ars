package main

import (
	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/lib4go/scheduler"
	"github.com/arsgo/lib4go/utility"
)

//BindLocalJobs 绑定本地JOB
func (a *AppServer) BindLocalJobs(jobs []cluster.JobItem) {
	scheduler.Stop()
	if jobs == nil || len(jobs) == 0 {
		a.Log.Info("未启用本地JOB,未配置或配置已删除")
		return
	}
	currentJobs := 0
	for _, v := range jobs {
		if !v.Enable {
			continue
		}
		er := a.scriptPool.PreLoad(v.Script, v.MinSize, v.MaxSize)
		if er != nil {
			a.Log.Errorf("JOB(%s)脚本(%s)加载失败:%v", v.Name, v.Script, er)
			continue
		}
		currentJobs++
		a.Log.Infof("::启动本地JOB(%s)[%s]", v.Name, v.Script)
		scheduler.AddTask(v.Trigger, scheduler.NewTask(v, func(job interface{}) {
			defer a.recover()
			item := job.(cluster.JobItem)
			context := base.NewInvokeContext(a.loggerName, utility.GetSessionID(), "{}", item.Params, "")
			context.Log.Infof(" -> 运行JOB(%s)[%s]", item.Name, item.Script)
			_, _, err := a.scriptPool.Call(item.Script, context)
			if err != nil {
				context.Log.Errorf("JOB运行异常:(%s)[%s],%v", item.Name, item.Script, err)
			} else {
				context.Log.Infof("JOB执行成功:%s", item.Name)
			}
		}))
	}
	if currentJobs > 0 {
		scheduler.Start()
	}
	a.Log.Infof("::local job count:%d ", currentJobs)
}
