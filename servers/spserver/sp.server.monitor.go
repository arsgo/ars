package main

import "time"

func (sp *SPServer) startMonitor() {
	go func() {
		tk := time.NewTicker(time.Second * 60)
		for {
			select {
			case <-tk.C:
				if sp.rpcClient.GetServiceCount() == 0 {
					sp.timerReloadRCServer.Push("rc server count is 0")
				}
			}
		}
	}()
START:
	if sp.clusterClient.WaitForConnected() {
		sp.Log.Info(" -> 已重新连接，重新发布服务")
		sp.resetCluster()
		goto START
	}
}
func (sp *SPServer) reloadRCServer(p ...interface{}) {
	items, err := sp.clusterClient.GetAllRCServers()
	sp.BindRCServer(items, err)
}
func (sp *SPServer) collectReporter(success int, failed int, err int) {
	if err > 0 {
		sp.timerReloadRCServer.Push("not find rc server")
	}
}
