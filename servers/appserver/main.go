package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/arsgo/ars/servers/config"
	"github.com/arsgo/lib4go/forever"
)
// go tool pprof appserver /tmp/profile175976149/mem.pprof

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	conf, err := config.Get()
	if err != nil {
		fmt.Println("获取配置文件失败:", err)
	}
	appServer, err := NewAPPServer(conf)
	if err != nil {
		os.Exit(100)
		return
	}
	f := forever.NewForever(appServer, appServer.Log, "appserver", "appserver")
	f.Start()
}
