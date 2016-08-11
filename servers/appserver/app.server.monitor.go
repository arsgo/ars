package main

import "time"

func (a *AppServer) startMonitor() {
	go func() {
		tk := time.NewTicker(time.Second * 60)
		for {
			select {
			case <-tk.C:
				if !a.disableRPC && a.rpcClient.GetServiceCount() == 0 {
					a.timerReloadRCServer.Push("可用 rc server 数为: 0")
				}
			}
		}
	}()
START:
	if a.clusterClient.WaitForConnected() {
		a.Log.Info(" |-> 已重新连接，重新发布服务")
		a.resetAppServer()
		goto START
	}
}

func (a *AppServer) reloadRCServer(p ...interface{}) {
	items, err := a.clusterClient.GetAllRCServers()
	a.BindRCServer(items, err)
}
func (a *AppServer) collectReporter(success int, failed int, err int) {
	if err > 0 {
		a.Log.Info(">>collectReporter")
		a.timerReloadRCServer.Push("可用 rc server 数为: 0")
	}
}