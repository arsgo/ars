package main

import (
	"fmt"
	"runtime"

	"github.com/arsgo/ars/servers/config"
	"github.com/arsgo/lib4go/forever"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	conf, err := config.Get()
	if err != nil {
		fmt.Println("获取配置文件失败:", err)
	}
	spServer, err := NewSPServer(conf)
	if err != nil {
		return
	}
	f := forever.NewForever(spServer, spServer.Log, "spserver", "spserver")
	f.Start()
}
