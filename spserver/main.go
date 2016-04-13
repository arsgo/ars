package main

import (
	"log"
	"runtime"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/forever"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fv := forever.NewForever("spserver", "spserver")
	result, err := fv.Manage(func() forever.IClose {
		spserver := cluster.NewSPServer()
		spserver.StartRPC()
		spserver.WatchServiceConfigChange()
		return spserver
	}, func(o forever.IClose) {
		o.Close()
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)

}
