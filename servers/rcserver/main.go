package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/colinyl/lib4go/forever"
)

func main() {
	fp, _ := os.Getwd()
	fmt.Println(fp)
	runtime.GOMAXPROCS(runtime.NumCPU())
	rcServer := NewRCServer()
	f := forever.NewForever(rcServer, rcServer.Log, "rcserver", "rcserver")
	f.Start()
}
