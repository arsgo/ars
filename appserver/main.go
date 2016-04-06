package main

import (
	"runtime"
	"time"

	"github.com/colinyl/ars/cluster"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	appServer := cluster.NewAPPServer()

	appServer.WatchRCServerChange(func(config []*cluster.RCServerConfig, err error) {
		appServer.BindRCServer(config, err)
	})

	appServer.WatchConfigChange(func(config *cluster.AppConfig, err error) error {
		return appServer.BindTask(config, err)
	})

	time.Sleep(time.Hour)
}
