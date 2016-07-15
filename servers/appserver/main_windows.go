package main

import (
	"runtime"

	"github.com/colinyl/lib4go/forever"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	appServer := NewAPPServer()
	f := forever.NewForever(appServer, appServer.Log, "appserver", "appserver")
	f.Start()
}
