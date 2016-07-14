package main

import (
	"os"
	"runtime"

	"github.com/colinyl/lib4go/forever"
)

func main() {

	fx, _ := os.OpenFile("C:\\tmp\\11.txt", os.O_WRONLY|os.O_CREATE|os.O_SYNC,
		0755)
	os.Stdout = fx
	os.Stderr = fx

	runtime.GOMAXPROCS(runtime.NumCPU())
	spServer := NewSPServer()
	f := forever.NewForever(spServer, spServer.Log, "spserver", "spserver")
	f.Start()
}
