package main

import (
	"runtime"

	"github.com/arsgo/lib4go/forever"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	spServer := NewSPServer()
	f := forever.NewForever(spServer, spServer.Log, "spserver", "spserver")
	f.Start()
}
