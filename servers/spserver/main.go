package main

import (
	"runtime"

	"github.com/colinyl/lib4go/forever"
)

func main() {

	//defer profile.Start(profile.CPUProfile).Stop()
	//defer profile.Start(profile.BlockProfile).Stop()
	//defer profile.Start(profile.MemProfile).Stop()

	runtime.GOMAXPROCS(runtime.NumCPU())
	spServer := NewSPServer()
	f := forever.NewForever(spServer, spServer.Log, "spserver", "spserver")
	f.Start()

}
