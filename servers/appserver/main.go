package main

import (
	"os"
	"runtime"

	"github.com/arsgo/lib4go/forever"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	appServer, err := NewAPPServer()
	if err != nil {
		os.Exit(100)
		return
	}
	f := forever.NewForever(appServer, appServer.Log, "appserver", "appserver")
	f.Start()
}
