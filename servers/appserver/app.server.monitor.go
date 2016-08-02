package main

import "time"

func (a *AppServer) startMonitor() {
	go func() {
		tk := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-tk.C:
				if a.rpcClient.GetServiceCount() == 0 {
					items, err := a.clusterClient.GetAllRCServerValues()
					if len(items) > 0 {
						a.Log.Info(" |-> 重新绑定rc server")
						a.BindRCServer(items, err)
						a.resetAppServer()
					}
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
