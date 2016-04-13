package cluster

import (
	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/ars/scheduler"
)

func (r *rcServer) BindScheduler(config *JobConfigs, err error) {
	if err != nil {
		r.Log.Error(err)
		return
	}
	scheduler.Stop()
	r.Log.Infof("job config has changed:%d", len(config.Jobs))
	if len(config.Jobs) == 0 {
		return
	}
	for _, v := range config.Jobs {
		scheduler.AddTask(v.Trigger, scheduler.NewTask(v.Name, func(name string) {
			consumers := r.getJobConsumers(name)
			if err != nil {
				r.Log.Infof("job [%s] download consumer error", name)
				r.Log.Error(err)
				return
			}
			total := config.Jobs[name].Concurrency
			index := 0
			for i := 0; i <= len(consumers); i++ {
				client := rpcservice.NewRPCClient(consumers[i])
				if client.Open() != nil {
					r.Log.Infof("open rpc port(%s) error ", consumers[i])
					continue
				}
				result, err := client.Request(name, "{}")
				client.Close()
				if err != nil {
					continue
				}
				if !ResultIsSuccess(result) {
					r.Log.Infof("job(%s) rpc (%s) error ", name, consumers[i])
					continue
				}
				index++
				if index >= total {
					continue
				}
			}
			if index >= total {
				r.Log.Infof("job %s executing success", name)
			} else {
				r.Log.Infof("job %s executing failed(%d/%d)", name, index, total)
			}

		}))
	}
	scheduler.Start()

}
