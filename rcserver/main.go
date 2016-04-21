package main

import (
	"runtime"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/forever"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rcServer := cluster.NewRCServer()
	f := forever.NewForever(rcServer, rcServer.Log, "rcserver", "rcserver")
	f.Start()
}
