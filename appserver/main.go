package main

import (
	"runtime"
	"sync"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/forever"
)

type service struct {
	svs   map[string]string
	mutex sync.Mutex
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	appServer := cluster.NewAPPServer()
	f := forever.NewForever(appServer, appServer.Log, "appserver", "appserver")
	f.Start()
}
