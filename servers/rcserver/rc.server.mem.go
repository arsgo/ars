package main

import (
	"runtime/debug"
	"time"

	"github.com/arsgo/lib4go/sysinfo"
)

func (rc *RCServer) clearMem() {
	tk := time.NewTicker(time.Second * 120)
	for {
		select {
		case <-tk.C:
			rc.snapLogger.Infof(" -> 清理内存...%dM", sysinfo.GetAPPMemory())
			debug.FreeOSMemory()
		}
	}
}
