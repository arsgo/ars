package main

import (
	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/rpcproxy"
	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/scheduler"
)

//BindJobScheduler 绑定RC服务器的JOB任务
func (rc *RCServer) BindJobScheduler(jobs map[string]cluster.TaskItem, err error) {
	if err != nil {
		rc.Log.Error(err)
		return
	}

	scheduler.Stop()
	rc.Log.Infof("job config has changed:%d", len(jobs))
	if len(jobs) == 0 {
		return
	}
	var jobCount int
	for _, v := range jobs {
		if v.Concurrency <= 0 {
			continue
		}
		jobCount++
		scheduler.AddTask(v.Trigger, scheduler.NewTask(v, func(v interface{}) {
			task := v.(cluster.TaskItem)
			consumers := rc.clusterClient.GetJobConsumers(task.Name)
			total := jobs[task.Name].Concurrency
			index := 0
			for i := 0; i < len(consumers); i++ {
				client := rpcservice.NewRPCClient(consumers[i])
				if client.Open() != nil {
					rc.Log.Infof("open rpc server(%s) error ", consumers[i])
					continue
				}
				result, err := client.Request(task.Name, "{}")
				client.Close()
				if err != nil {
					rc.Log.Error(err)
					continue
				}
				if !rpcproxy.ResultIsSuccess(result) {
					rc.Log.Infof("call job(%s - %v) failed %s", task.Name, consumers[i], result)
					continue
				} else {
					rc.Log.Infof("call job(%s - %s) success", task.Name, consumers[i])
				}
				index++
				if index >= total {
					continue
				}
			}
			rc.Log.Infof("job(%s) has executed (%d/%d),consumers:%d", task.Name, index, total, len(consumers))

		}))
	}
	if jobCount > 0 {
		scheduler.Start()
	}

}
