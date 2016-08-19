package main

import (
	"strings"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/lib4go/scheduler"
	"github.com/arsgo/lib4go/utility"
)

//bindLocalJobs 绑定本地JOB
func (a *AppServer) bindLocalJobs(tasks []cluster.TaskItem) {
	jobs := a.getJobConsumerTask(tasks)
	scheduler.Stop()
	if jobs == nil || len(jobs) == 0 {
		a.Log.Info("--> 未启用本地JOB,未配置或配置已删除")
		return
	}
	currentJobs := 0
	for _, v := range jobs {
		if v.Disable {
			continue
		}
		if strings.EqualFold(v.Trigger, "") {
			a.Log.Errorf("JOB(%s)未配置trigger", v.Name)
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
			context.Log.Infof("--> 运行JOB(%s)", item.Name)
			results, _, err := a.scriptPool.Call(item.Script, context)
			if err != nil || len(results) != 1 || !base.ResultIsSuccess(results[0]) {
				context.Log.Errorf("--> JOB运行异常:(%s)[%s],result:%v,err:%v", item.Name, item.Script, results, err)
				return
			}
			context.Log.Infof("--> JOB执行成功:%s", item.Name)

		}))
	}
	if currentJobs > 0 {
		scheduler.Start()
	}
	a.Log.Infof("::local job count:%d ", currentJobs)
}
