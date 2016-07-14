package main

import (
	"os"
	"runtime"
	"syscall"

	"github.com/colinyl/lib4go/forever"
)

func main() {
	logFile, _ := os.OpenFile("/ext/rcserver", os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
	syscall.Dup2(int(logFile.Fd()), 1)
	syscall.Dup2(int(logFile.Fd()), 2)

	runtime.GOMAXPROCS(runtime.NumCPU())
	rcServer := NewRCServer()
	f := forever.NewForever(rcServer, rcServer.Log, "rcserver", "rcserver")
	f.Start()
}
