package main

import (
	"fmt"
	"runtime"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/forever"
	"github.com/colinyl/lib4go/utility"
)

func main() {

	fmt.Println(utility.GetParams("name=colin&sex=100&sex=200&order=1234567890"))

	runtime.GOMAXPROCS(runtime.NumCPU())
	appServer := cluster.NewAPPServer()
	f := forever.NewForever(appServer, appServer.Log, "appserver", "appserver")
	f.Start()

}
