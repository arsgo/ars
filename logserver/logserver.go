package main

import (
	"runtime"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/forever"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	spServer := cluster.NewSPServer()
	f := forever.NewForever(spServer, spServer.Log, "spserver", "spserver")
	f.Start()

}
