package main

import (
	"runtime/debug"
	"time"

	"github.com/arsgo/lib4go/sysinfo"
)

func (a *AppServer) clearMem() {
	tk := time.NewTicker(time.Second * 120)
	for {
		select {
		case <-tk.C:
			a.snapLogger.Infof(" -> 清理内存...%dM", sysinfo.GetAPPMemory())
			debug.FreeOSMemory()
		}
	}
}
