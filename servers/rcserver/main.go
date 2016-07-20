package main

import (
	"runtime"

	"github.com/colinyl/lib4go/forever"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rcServer := NewRCServer()
	f := forever.NewForever(rcServer, rcServer.Log, "rcserver", "rcserver")
	f.Start()
}
