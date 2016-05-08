package main

import (
	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/rpcclient"
	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/scheduler"
)

//BindJobScheduler 绑定RC服务器的JOB任务
func (rc *RCServer) BindJobScheduler(config *cluster.JobItems, err error) {
	if err != nil {
		rc.Log.Error(err)
		return
	}

	scheduler.Stop()
	rc.Log.Infof("job config has changed:%d", len(config.Jobs))
	if len(config.Jobs) == 0 {
		return
	}
	var jobCount int
	for _, v := range config.Jobs {
		if v.Concurrency <= 0 {
			continue
		}
		jobCount++
		scheduler.AddTask(v.Trigger, scheduler.NewTask(v.Name, func(v interface{}) {
			name := v.(string)
			consumers := rc.clusterClient.GetJobConsumers(name)
			if err != nil {
				rc.Log.Infof("job ", name, " download consumers error", err)
				return
			}
			total := config.Jobs[name].Concurrency
			index := 0
			for i := 0; i < len(consumers); i++ {
				client := rpcservice.NewRPCClient(consumers[i])
				if client.Open() != nil {
					rc.Log.Infof("open rpc server(%s) error ", consumers[i])
					continue
				}
				result, err := client.Request(name, "{}")
				client.Close()
				if err != nil {
					rc.Log.Error(err)
					continue
				}
				if !rpcclient.ResultIsSuccess(result) {
					rc.Log.Infof("call job(%s - %s) failed ", name, consumers[i])
					continue
				} else {
					rc.Log.Infof("call job(%s - %s) success", name, consumers[i])
				}
				index++
				if index >= total {
					continue
				}
			}
			rc.Log.Infof("job(%s) has executed (%d/%d),consumers:%d", name, index, total, len(consumers))

		}))
	}
	if jobCount > 0 {
		scheduler.Start()
	}

}
