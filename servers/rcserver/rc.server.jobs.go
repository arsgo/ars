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

func (rc *RCServer) getJobs(jobs []cluster.JobItem) (jobMap map[string]cluster.JobItem) {
	jobMap = make(map[string]cluster.JobItem)
	if jobs == nil {
		return
	}
	for _, v := range jobs {
		jobMap[v.Name] = v
	}
	return
}

//BindJobScheduler 绑定RC服务器的JOB任务
func (rc *RCServer) BindJobScheduler(jobItems []cluster.JobItem, err error) {
	scheduler.Stop()
	if err != nil || len(jobItems) == 0 {
		rc.Log.Info(" -> 未配置job或已删除")
		return
	}
	jobs := rc.getJobs(jobItems)
	var currentJobs int
	for _, v := range jobs {
		if v.Concurrency <= 0 || v.Disable {

			continue
		}
		currentJobs++
		rc.Log.Infof("::start job: %s", v.Name)
		scheduler.AddTask(v.Trigger, scheduler.NewTask(v, func(v interface{}) {
			task := v.(cluster.JobItem)
			consumers := rc.clusterClient.GetJobConsumers(task.Name)
			rc.startJobFlow(task, consumers, utility.GetSessionID())
		}))
	}
	if currentJobs > 0 {
		scheduler.Start()
	}
	//	rc.Log.Infof("::当前已启动的job数:%d", currentJobs)
}

func (rc *RCServer) startJobFlow(task cluster.JobItem, consumers []string, session string) {
	if len(consumers) == 0 {
		rc.Log.Errorf("--> 无法运行 job (%s),无可用的consumer", task.Name)
		rc.schedulerCollector.Error(task.Name)
		return
	}
	avaliable := base.NewAvaliableMap(consumers)
	data := make(chan int, task.Concurrency)
	clogger, err := logger.NewSession(rc.loggerName, session)
	if err != nil {
		clogger = rc.Log
	}
	clogger.Infof("--> 运行 job[ %s %d,%d]", task.Name, task.Concurrency, len(consumers))
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
		rc.schedulerCollector.Failed(task.Name)
		clogger.Errorf("--> end job(%s)未全部执行成功[%d/%d],consumer:%d", task.Name, success, task.Concurrency, len(consumers))
	} else {
		rc.schedulerCollector.Success(task.Name)
		clogger.Infof("--> end job(%s)全部执行成功(%d次)", task.Name, task.Concurrency)
	}
}

func (rc *RCServer) runJob(task cluster.JobItem, consumer string, session string) (err error) {
	defer func() {
		err = rc.recover()
	}()
	client := rpcservice.NewRPCClient(consumer, rc.loggerName)
	if client.Open() != nil {
		err = fmt.Errorf("执行失败job(%s)失败,无法连接到服务器: %s", task.Name, consumer)
		return
	}
	result, err := client.Request(task.Name, "{}", session)
	client.Close()
	if err != nil {
		err = fmt.Errorf("调用job(%s)失败,%v", task.Name, err)
		return
	}
	if !base.ResultIsSuccess(result) {
		err = fmt.Errorf("--> 调用job(%s-%s)返回失败:%s", task.Name, consumer, result)
	}
	return
}
