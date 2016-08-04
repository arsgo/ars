package main

import "time"

func (sp *SPServer) startMonitor() {
	go func() {
		tk := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-tk.C:
				if sp.rpcClient.GetServiceCount() == 0 {
					items, err := sp.clusterClient.GetAllRCServers()
					if len(items) > 0 {
						sp.Log.Info(" |-> 重新绑定rc server")
						sp.BindRCServer(items, err)
						sp.resetCluster()
					}
				}
			}
		}
	}()
START:
	if sp.clusterClient.WaitForConnected() {
		sp.Log.Info(" |-> 已重新连接，重新发布服务")
		sp.resetCluster()
		goto START
	}
}
