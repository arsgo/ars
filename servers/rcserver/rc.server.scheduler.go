package main

import (
	"fmt"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/base/rpcservice"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/lib4go/logger"
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
			rc.startJobFlow(task, consumers, utility.GetSessionID())
		}))
	}
	if currentJobs > 0 {
		scheduler.Start()
	}
	rc.Log.Infof("::当前已启动的scheduler数:%d", currentJobs)
}

func (rc *RCServer) startJobFlow(task cluster.JobItem, consumers []string, session string) {
	avaliable := base.NewAvaliableMap(consumers)
	data := make(chan int, task.Concurrency)
	clogger, err := logger.NewSession(rc.loggerName, session)
	if err != nil {
		clogger = rc.Log
	}
	clogger.Infof(" -> 运行scheduler(%s)[%d,%d]", task.Name, task.Concurrency, len(consumers))
	for i := 0; i < task.Concurrency; i++ {
		data <- i
	}
	success := 0
START:
	for {
		select {
		case <-data:
			consumer, err := avaliable.Get()
			if err != nil {
				break START
			}
			go func(consumer string) {
				if r := rc.runJob(task, consumer, session); r != nil {
					avaliable.Remove(consumer)
					clogger.Error(r)
					data <- 0
				} else {
					success++
				}
			}(consumer)
		default:
			if success == task.Concurrency {
				break START
			}
		}
	}
	if success < task.Concurrency {
		clogger.Errorf(" -> scheduler(%s)未完全执行成功,已执行: %d次, 总共: %d次", task.Name, success, task.Concurrency)
	} else {
		clogger.Infof(" -> 执行成功scheduler(%s)", task.Name)
	}
}

func (rc *RCServer) runJob(task cluster.JobItem, consumer string, session string) (err error) {
	defer func() {
		err = rc.recover()
	}()
	client := rpcservice.NewRPCClient(consumer, rc.loggerName)
	if client.Open() != nil {
		err = fmt.Errorf("执行失败scheduler(%s)失败,无法连接到服务器: %s", task.Name, consumer)
		return
	}
	result, err := client.Request(task.Name, "{}", session)
	client.Close()
	if err != nil {
		err = fmt.Errorf("调用scheduler(%s)失败,%v", task.Name, err)
		return
	}
	if !base.ResultIsSuccess(result) {
		err = fmt.Errorf(" -> 调用scheduler(%s-%s)返回失败:%s", task.Name, consumer, result)
	}
	return
}
