package main

import (
	"runtime"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/forever"
	"github.com/colinyl/lib4go/logger"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	appServer := cluster.NewAPPServer()
	defer func() {
		if err := recover(); err != nil {
			log, _ := logger.New("exp", true)
			log.Error(err)
		}
	}()
	f := forever.NewForever(appServer, appServer.Log, "appserver", "appserver")
	f.Start()

}
