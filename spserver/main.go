package main

import (
	"runtime"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/forever"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	spserver := cluster.NewSPServer()
	fv := forever.NewForever("spserver", "spserver")
	result, err := fv.Manage(func() {
		spserver.StartRPC()
		spserver.WatchServiceConfigChange()
	}, func() {
		spserver.Close()
	})
	if err != nil {
		spserver.Log.Error(err)
		return
	}
	spserver.Log.Info(result)

}
