package main

import (
	"runtime"

	"github.com/colinyl/lib4go/forever"
)

func main() {
	/*	logFile, _ := os.OpenFile("./appserver.dup", os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
		syscall.Dup2(int(logFile.Fd()), 1)
		syscall.Dup2(int(logFile.Fd()), 2) *
	*/
	runtime.GOMAXPROCS(runtime.NumCPU())
	appServer := NewAPPServer()
	f := forever.NewForever(appServer, appServer.Log, "appserver", "appserver")
	f.Start()
}
