package main

import (
	"fmt"
	"runtime"

	"github.com/colinyl/lib4go/forever"
	"github.com/colinyl/lib4go/utility"
)

func main() {
	fmt.Println(utility.GetExcPath("./ars.conf.json"))
	runtime.GOMAXPROCS(runtime.NumCPU())
	rcServer := NewRCServer()
	f := forever.NewForever(rcServer, rcServer.Log, "rcserver", "rcserver")
	f.Start()
}
