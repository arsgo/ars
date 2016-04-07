package main

import (
	"runtime"
	"time"

	"github.com/colinyl/ars/cluster"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	spserver := cluster.NewSPServer()
	spserver.StartRPC()
	spserver.WatchServiceConfigChange()

	time.Sleep(time.Hour)

}
