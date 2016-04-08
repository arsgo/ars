package main

import (
	"runtime"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/forever"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rcServer := cluster.NewRCServer()
	fv := forever.NewForever("rcserver", "rcserver")
	result, err := fv.Manage(func(){
		rcServer.Bind()
		rcServer.StartRPCServer()
		rcServer.WatchJobChange(func(config *cluster.JobConfigs, err error) {

		})
		rcServer.WatchServiceChange(func(services map[string][]string, err error) {
			rcServer.BindSPServer(services)
		})
		rcServer.StartSnapValue()
	}, func() {

	})

	if err != nil {
		rcServer.Log.Error(err)
		return
	}
	rcServer.Log.Info(result)

}
