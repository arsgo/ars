package main

import (
	"os"
	"runtime"

	"github.com/colinyl/lib4go/forever"
)

func main() {
	logFile, _ := os.OpenFile("/ext/spserver", os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
	syscall.Dup2(int(logFile.Fd()), 1)
	syscall.Dup2(int(logFile.Fd()), 2)

	runtime.GOMAXPROCS(runtime.NumCPU())
	spServer := NewSPServer()
	f := forever.NewForever(spServer, spServer.Log, "spserver", "spserver")
	f.Start()
}
