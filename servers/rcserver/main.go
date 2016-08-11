package main

import (
	"os"
	"runtime"

	"github.com/arsgo/lib4go/forever"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rcServer, err := NewRCServer()
	if err != nil {
		os.Exit(100)
		return
	}
	f := forever.NewForever(rcServer, rcServer.Log, "rcserver", "rcserver")
	f.Start()
}
