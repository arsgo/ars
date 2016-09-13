package main

import (
	"strings"

	"github.com/arsgo/ars/cluster"
)

//BindCrossAccess 绑定垮域访问服务
func (rc *RCServer) BindCrossAccess(task cluster.RCServerTask) (err error) {
	if len(task.CrossDomainAccess) == 0 {
		rc.Log.Info(" -> 未配置跨越访问或已删除")
	}
	rc.ResetCrossDomainServices(task)
	rc.WatchCrossDomain(task)
	rc.PublishNow()
	return
}

func (rc *RCServer) checkDomain(domain string, item cluster.CrossDoaminAccessItem) bool {
	if !strings.EqualFold("@"+domain, rc.clusterClient.GetDomainName()) {
		return true
	}
	return !rc.isInServerList(item.Servers)
}

//ResetCrossDomainServices 重置跨域服务
func (rc *RCServer) ResetCrossDomainServices(task cluster.RCServerTask) {

	//添加、关闭、更新跨域服务器
	localCrossServices := rc.crossServices.GetAll()
	//添加不存在的域和服务器
	for domain, item := range task.CrossDomainAccess {
		if item.Disable {
			continue
		}
		if _, ok := localCrossServices[domain]; ok {
			continue
		}
		if !rc.checkDomain(domain, item) {
			rc.Log.Error(" -> 域配置有误，不能与当前域相同:", domain)
			continue
		}
		crossData := item.GetServicesMap(domain) //转换为服务映射表
		rc.crossServices.Set(domain, crossData)  //添加到服务列表
	}
	//删除，更新服务
	for domain, svs := range localCrossServices {
		if _, ok := task.CrossDomainAccess[domain]; !ok {
			rc.crossServices.Delete(domain) //不存在域,则删除
			continue
		}
		if v, ok := task.CrossDomainAccess[domain]; ok && v.Disable {
			rc.crossServices.Delete(domain) //域已禁用,则删除
			rc.Log.Error(" -> 域配置已禁用:", domain)
			continue
		}

		//检查本地服务是否与远程服务一致
		currentServices := svs.(cluster.RPCServices)                            //本地服务
		remoteServices := task.CrossDomainAccess[domain].GetServicesMap(domain) //远程服务
		//删除更新服务
		for name := range currentServices {
			if _, ok := remoteServices[name]; !ok {
				delete(currentServices, name) //远程不存在，则删除本地服务
			} else {
				currentServices[name] = task.CrossDomainAccess[domain].Servers //覆盖本地服务
			}
		}
		//添加服务
		for name := range remoteServices {
			if _, ok := currentServices[name]; !ok {
				currentServices[name] = task.CrossDomainAccess[domain].Servers //本地不存在则添加服务
			}
		}
	}
}

//WatchCrossDomain 监控外部域
func (rc *RCServer) WatchCrossDomain(task cluster.RCServerTask) {
	//关闭域
	localDomains := rc.crossDomain.GetAll()
	for domain, clt := range localDomains {
		if v, ok := task.CrossDomainAccess[domain]; !ok || (ok && v.Disable) {
			rc.Log.Info(" -> 关闭外部域:", domain)
			client := clt.(cluster.IClusterClient)
			client.Close()
			rc.crossDomain.Delete(domain)
		}
	}

	//监控域
	for domain, v := range task.CrossDomainAccess {
		if !rc.checkDomain(domain, v) || v.Disable {
			continue
		}
		//为cluster类型时,添加监控
		if _, ok := rc.crossDomain.Get(domain); !ok {
			var clusterClient cluster.IClusterClient
			var err error
			if rc.isInServerList(v.Servers) {
				rc.Log.Info(" -> 启动外部域:", domain)
				clusterClient, err = cluster.NewDomainClusterClientHandler(domain, rc.conf.IP, rc.loggerName, rc.clusterClient.GetHandler())
			} else {
				rc.Log.Info(" -> 启动外部域:", domain)
				clusterClient, err = cluster.NewDomainClusterClient(domain, rc.conf.IP, rc.loggerName, v.Servers...)
			}
			if err != nil {
				rc.Log.Error(err)
				continue
			}

			//将服务添加到服务列表
			rc.crossDomain.Set(domain, clusterClient)
			rc.crossServices.Set(domain, v.GetServicesMap(domain))

			//监控外部RC服务器变化,变化后更新本地服务
			go func(domain string) {
				if !rc.IsMaster {
					return
				}
				defer rc.recover()
				clusterClient.WatchRCServerChange(func(items []*cluster.RCServerItem, err error) {
					if !rc.crossDomain.Exists(domain) {
						return
					}
					rc.bindCrossServices(domain, items)
					rc.PublishNow()
				})
			}(domain)
		}
	}
}
func (rc *RCServer) getDomainIPs(items []*cluster.RCServerItem) []string {
	var ips = []string{}
	for _, v := range items {
		ips = append(ips, v.Address)
	}
	return ips
}
func (rc *RCServer) bindCrossServices(domain string, items []*cluster.RCServerItem) {
	ips := rc.getDomainIPs(items)
	services, ok := rc.crossServices.Get(domain)
	if !ok {
		return
	}
	allServices := services.(cluster.RPCServices)
	for name := range allServices {
		allServices[name] = ips
	}
	rc.crossServices.Set(domain, allServices)
}
func (rc *RCServer) isInServerList(lst []string) bool {
	for _, i := range lst {
		exits := false
		for _, j := range rc.clusterServers {
			if strings.EqualFold(i, j) {
				exits = true
				break
			}
		}
		if !exits {
			return false
		}
	}
	return true
}
