package main

import (
	"runtime"

	"github.com/colinyl/ars/cluster"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rcServer := cluster.NewRCServer()
	rcServer.Bind()
	rcServer.StartRPCServer()
	rcServer.WatchJobChange(func(config *cluster.JobConfigs, err error) {	

	})
	rcServer.WatchServiceChange(func(services map[string][]string, err error) {
		rcServer.BindSPServer(services)
	})
	rcServer.StartSnapValue()

}
