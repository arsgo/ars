package main

import (
	"runtime"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/forever"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	appServer := cluster.NewAPPServer()
	f := forever.NewForever(appServer, appServer.Log, "appserver", "appserver")
	f.Start()
}
