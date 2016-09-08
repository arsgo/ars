package main

import "time"

func (a *AppServer) clearMem() {
	tk := time.NewTicker(time.Second * 201)
	for {
		select {
		case <-tk.C:
			//	a.snapLogger.Infof(" -> 清理内存...%dM", sysinfo.GetAPPMemory())
			//debug.FreeOSMemory()
		}
	}
}
