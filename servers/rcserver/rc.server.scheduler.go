package main

import (
	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/base/rpcservice"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/lib4go/scheduler"
	"github.com/arsgo/lib4go/utility"
)

//BindJobScheduler 绑定RC服务器的JOB任务
func (rc *RCServer) BindJobScheduler(jobs map[string]cluster.JobItem, err error) {
	scheduler.Stop()
	if err != nil || len(jobs) == 0 {
		rc.Log.Info("获取scheduler配置失败或未配置")
		return
	}

	var currentJobs int
	for _, v := range jobs {
		if v.Concurrency <= 0 || !v.Enable {
			continue
		}
		currentJobs++
		rc.Log.Infof("::启动scheduler: %s", v.Name)
		scheduler.AddTask(v.Trigger, scheduler.NewTask(v, func(v interface{}) {
			task := v.(cluster.JobItem)
			consumers := rc.clusterClient.GetJobConsumers(task.Name)
			total := jobs[task.Name].Concurrency
			runSucess := 0
			for i := 0; i < len(consumers); i++ {
				rc.Log.Infof(" -> 运行scheduler(%s)[%s]", task.Name, consumers[i])
				client := rpcservice.NewRPCClient(consumers[i], rc.loggerName)
				if client.Open() != nil {
					rc.Log.Errorf("运行scheduler失败,无法连接到服务器: %s", consumers[i])
					continue
				}
				result, err := client.Request(task.Name, "{}", utility.GetSessionID())
				client.Close()
				if err != nil {
					rc.Log.Errorf("调用scheduler失败,%v", err)
					continue
				}
				if !base.ResultIsSuccess(result) {
					rc.Log.Errorf(" -> 调用scheduler(%s-%s)返回失败:%s", task.Name, consumers[i], result)
					continue
				} else {
					rc.Log.Infof(" -> scheduler(%s-%s)执行成功", task.Name, consumers[i])
				}
				runSucess++
				if runSucess >= total {
					continue
				}
			}
			if runSucess < total {
				rc.Log.Errorf(" -> scheduler(%s)未完全执行成功,已执行: %d次, 总共: %d次", task.Name, runSucess, total)
			}
		}))
	}
	if currentJobs > 0 {
		scheduler.Start()
	}
	rc.Log.Infof("::当前已启动的scheduler数:%d", currentJobs)
}
