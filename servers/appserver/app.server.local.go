package main

import (
	"fmt"
	"strings"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/lib4go/scheduler"
	"github.com/arsgo/lib4go/utility"
)

//bindLocalJobs 绑定本地JOB
func (a *AppServer) bindLocalJobs(tasks []cluster.TaskItem) {
	jobs := a.getJobLocalsTask(tasks)
	scheduler.Stop()
	if jobs == nil || len(jobs) == 0 {
		a.OnLocalJobCloseAll()
		return
	}
	currentJobs := 0
	for _, v := range jobs {
		if v.Disable {
			a.OnLocalJobClose(v)
			continue
		}
		if strings.EqualFold(v.Trigger, "") {
			a.Log.Errorf(" -> JOB(%s)未配置trigger", v.Name)
			a.jobLocalCollector.Error(v.Name)
			a.OnLocalJobClose(v)
			continue
		}
		er := a.scriptPool.PreLoad(v.Script, v.MinSize, v.MaxSize)
		if er != nil {
			a.Log.Errorf("JOB(%s)脚本(%s)加载失败:%v", v.Name, v.Script, er)
			a.jobLocalCollector.Error(v.Name)
			a.OnLocalJobClose(v)
			continue
		}
		currentJobs++
		a.OnLocalJobCreate(v)
		scheduler.AddTask(v.Trigger, scheduler.NewTask(v, func(job interface{}) {
			defer a.recover()
			item := job.(cluster.TaskItem)
			context := base.NewInvokeContext(item.Name, base.TN_JOB_LOCAL, a.loggerName, utility.GetSessionID(), "{}", item.Params, "")
			context.Log.Infof("--> 运行JOB(%s)", item.Name)
			results, _, err := a.scriptPool.Call(item.Script, context)

			if err != nil || len(results) != 1 || (!base.ResultIsSuccess(results[0]) && !strings.EqualFold(fmt.Sprintf("%s", results[0]), "true")) {
				context.Log.Errorf("--> JOB运行异常:(%s)[%s],result:%v,err:%v", item.Name, item.Script, results, err)
				a.jobLocalCollector.Failed(item.Name)
				return
			}
			a.jobLocalCollector.Success(item.Name)
			context.Log.Infof("--> JOB执行成功:%s", item.Name)

		}))
	}
	if currentJobs > 0 {
		scheduler.Start()
	}
	//a.Log.Infof("::local job count:%d ", currentJobs)
}

//OnLocalJobCreate  创建本地JOB服务
func (a *AppServer) OnLocalJobCreate(task cluster.TaskItem) (path string) {
	path, err := a.clusterClient.CreateLocalJob(task.Name, a.snap.getJobLocalSnap())
	if err != nil {
		a.Log.Error("local job 创建失败: ", err)
		a.jobLocalCollector.Error(task.Name)
		return
	}
	a.localJobPaths.Set(task.Name, path)
	a.Log.Infof("::start local job: [%s] %s", task.Name, task.Script)
	return
}

//OnLocalJobClose close job
func (a *AppServer) OnLocalJobClose(task cluster.TaskItem) {
	p, ok := a.localJobPaths.Get(task.Name)
	if ok {
		path := p.(string)
		a.Log.Error(" -> 关闭本地job:", path)
		err := a.clusterClient.CloseNode(path)
		if err != nil {
			return
		}
		a.localJobPaths.Delete(task.Name)
	}
}

//OnLocalJobCloseAll 关闭本地JOB服务
func (a *AppServer) OnLocalJobCloseAll() {
	paths := a.localJobPaths.GetAllAndClear()
	if len(paths) == 0 {
		a.Log.Info("--> 本地JOB未配置或已删除")
		return
	}
	for name, p := range paths {
		path := p.(string)
		a.Log.Info(" -> 关闭本地job:", name)
		err := a.clusterClient.CloseNode(path)
		if err != nil {
			return
		}
		a.localJobPaths.Delete(name)
		return
	}

}
