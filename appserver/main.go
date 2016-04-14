package main

import (
	"log"
	"runtime"
	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/forever"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	fv := forever.NewForever("appserver", "appserver")
	result, err := fv.Manage(func() forever.IClose {
		appServer := cluster.NewAPPServer()
		appServer.WatchRCServerChange(func(config []*cluster.RCServerConfig, err error) {
			appServer.BindRCServer(config, err)
		})

		appServer.WatchConfigChange(func(config *cluster.AppConfig, err error) error {
			return appServer.BindTask(config, err)
		})
		return appServer
	}, func(o forever.IClose) {
		o.Close()
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)

}
