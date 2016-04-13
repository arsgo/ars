package main

import (
	"log"
	"runtime"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/forever"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fv := forever.NewForever("rcserver", "rcserver")
	result, err := fv.Manage(func() forever.IClose {
		rcServer := cluster.NewRCServer()
		rcServer.Bind()
		rcServer.StartRPCServer()
		rcServer.WatchJobChange(func(config *cluster.JobConfigs, err error) {
			rcServer.BindScheduler(config, err)
		})
		rcServer.WatchServiceChange(func(services map[string][]string, err error) {
			rcServer.BindSPServer(services)
		})
		rcServer.StartSnapValue()
		return rcServer
	}, func(o forever.IClose) {
		o.Close()
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)

}
