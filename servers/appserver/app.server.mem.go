package main

/*
func (a *AppServer) clearMem() {
	tk1 := time.NewTicker(time.Second * 10)
	tk := time.NewTicker(time.Second * 61)
	for {
		select {
		case <-tk1.C:
			a.snapLogger.Infof(" -> 内存...%dM", sysinfo.GetAPPMemory())
		case <-tk.C:
			a.snapLogger.Infof(" -> 清理内存...%dM", sysinfo.GetAPPMemory())
			//	debug.FreeOSMemory()
		}
	}
}
*/
