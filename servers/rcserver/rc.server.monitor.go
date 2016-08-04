package main

import (
	"strings"
	"time"

	"github.com/arsgo/ars/cluster"
)

//startMonitor 启动监控服务
func (rc *RCServer) startMonitor() {

	go func() {
		tk := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-tk.C:
				if rc.needBindRPCService() {
					rc.Log.Debug("~~~~~~rebind rpc service")
					rc.rebindLocalServices()
				}
			}
		}
	}()

START:
	if rc.clusterClient.WaitForConnected() {
		if rc.IsMaster {
			rc.Log.Info(" |-> 已重新连接，重新发布服务")
			rc.PublishAll()
		}
		goto START
	}
}
func (rc *RCServer) needBindRPCService() bool {
	nmap := make(map[string]bool)
	cdomain := strings.Replace(strings.TrimLeft(rc.domain, "/"), "/", ".", -1)
	nmap[cdomain] = false
	all := rc.crossDomain.GetAll()
	for i := range all {
		nmap[i] = false
	}
	services := rc.spRPCClient.GetServices()
	for sv := range services {
		index := strings.LastIndex(sv, "@")
		domain := sv[index+1:]
		if _, ok := nmap[domain]; ok {
			nmap[domain] = true
		}
	}
	for _, v := range nmap {
		if !v {
			return true
		}
	}
	return false
}

//rebindLocalServices 重新绑定本地服务
func (rc *RCServer) rebindLocalServices() (err error) {
	lst, err := rc.clusterClient.GetSPServerServices()
	if err != nil {
		rc.Log.Error(err)
		return
	}
	rc.currentServices.Set("*", lst)
	err = rc.resetCrossDomainServices()
	if err != nil {
		return
	}
	services := rc.MergeService()
	rc.BindSPServers(services, nil)
	return

}
func (rc *RCServer) resetCrossDomainServices() (err error) {
	task, err := rc.clusterClient.GetRCServerTask()
	if err != nil {
		rc.Log.Error(err)
		return
	}
	rc.ResetCrossDomainServices(task)
	rc.WatchCrossDomain(task)
	domains := rc.crossDomain.GetAll()
	rc.Log.Debug("domains:", domains)
	for domain, cst := range domains {
		serveritems, err := cst.(cluster.IClusterClient).GetAllRCServers()
		if err != nil {
			break
		}
		rc.bindCrossServices(domain, serveritems)
	}
	return
}
