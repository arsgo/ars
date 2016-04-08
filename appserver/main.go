package main

import (
	"runtime"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/forever"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	appServer := cluster.NewAPPServer()
	fv := forever.NewForever("appserver", "appserver")
	result, err := fv.Manage(func() {
		appServer.WatchRCServerChange(func(config []*cluster.RCServerConfig, err error) {
			appServer.BindRCServer(config, err)
		})

		appServer.WatchConfigChange(func(config *cluster.AppConfig, err error) error {
			return appServer.BindTask(config, err)
		})
	}, func() {

	})
	if err != nil {
		appServer.Log.Error(err)
		return
	}
	appServer.Log.Info(result)

}
