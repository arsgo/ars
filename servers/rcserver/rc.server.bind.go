package main

import (
	"sort"
	"strings"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/servers/config"
)

//BindRCServer 绑定服务
func (rc *RCServer) BindRCServer() (err error) {
	rc.snap.Address = rc.rcRPCServer.Address
	rc.snap.Path, err = rc.clusterClient.CreateRCServer(rc.snap.GetSnap())
	if err != nil {
		return
	}
	rc.clusterClient.WatchRCServerChange(func(items []*cluster.RCServerItem, err error) {
		rc.IsMaster = rc.IsMasterServer(items, rc.snap.Path)
		if rc.IsMaster {
			//as master
			rc.snap.Server = SERVER_MASTER
			rc.Log.Info("current server is ", SERVER_MASTER)

			rc.clusterClient.WatchJobConfigChange(func(config *cluster.JobItems, err error) {
				rc.BindJobScheduler(config, err)
			})
			rc.clusterClient.WatchServiceProviderChange(func() {
				rc.UpdateLocalService()
			})
			rc.clusterClient.WatchRCTaskChange(func(task cluster.RCServerTask, err error) {
				rc.BindCrossAccess(task)
			})

		} else {
			//as slave
			rc.Log.Info("current server is ", SERVER_SLAVE)
			rc.clusterClient.WatchRPCServiceChange(func(services map[string][]string, err error) {
				rc.spRPCClient.ResetRPCServer(services)
			})
		}

	})

	return
}

func (rc *RCServer) BindCrossAccess(task cluster.RCServerTask) (err error) {
	rc.crossLock.Lock()
	defer rc.crossLock.Unlock()

	//移除所有监控
	for domain, client := range rc.crossDomain {
		client.Close()
		delete(rc.crossDomain, domain)
	}

	//移除域和服务
	for domain, services := range rc.crossService {
		if _, ok := task.CrossDomainAccess[domain]; !ok {
			delete(rc.crossService, domain) //不存在域,则删除
		}
		for onesvs := range services {
			for _, v := range task.CrossDomainAccess[domain].Services {
				if strings.EqualFold(v, onesvs) {
					rc.crossService[domain][onesvs] = task.CrossDomainAccess[domain].Servers //存在服务,则更新IP列表
					continue
				}
			}
			delete(rc.crossService[domain], onesvs) //不存在服务,则移除服务
		}
	}

	for domain, v := range task.CrossDomainAccess {
		//为cluster类型时,添加监控
		if _, ok := rc.crossDomain[domain]; !ok && strings.EqualFold(strings.ToLower(v.Type), "cluster") {
			rc.crossDomain[domain], err = cluster.GetClusterClient(domain, config.Get().IP, v.Servers...)
			if err != nil {
				continue
			}
			rc.crossService[domain] = make(map[string][]string)
			for _, svs := range v.Services {
				rc.crossService[domain][svs] = v.Servers
			}
			//监控外部RC服务器变化,变化后更新本地服务
			go rc.crossDomain[domain].WatchRCServerChange(func(items []*cluster.RCServerItem, err error) {
				rc.crossLock.Lock()
				defer rc.crossLock.Unlock()
				var ips = []string{}
				for _, v := range items {
					ips = append(ips, v.Server)
				}
				for service := range rc.crossService[domain] {
					rc.crossService[domain][service] = ips
				}
			})
		} else if strings.EqualFold(strings.ToLower(v.Type), "proxy") {
			//为 proxy类型时,直接添加到服务列表
			for _, service := range v.Services {
				rc.crossService[domain][service] = v.Servers
			}

		}
	}
	//重新发布服务
	rc.clusterClient.PublishRPCServices(rc.crossService)
	rc.UpdateLocalService()
	return

}

//IsMasterServer 检查当前RC Server是否是Master
func (rc *RCServer) IsMasterServer(items []*cluster.RCServerItem, path string) bool {
	var servers []string
	for _, v := range items {
		servers = append(servers, v.Address)
	}

	sort.Sort(sort.StringSlice(servers))
	return len(servers) == 0 || strings.HasSuffix(path, servers[0])
}

//UpdateLocalService 更新本定服务列表
func (rc *RCServer) UpdateLocalService() {
	service, err := rc.clusterClient.GetRPCService()
	if err != nil {
		rc.spRPCClient.ResetRPCServer(service)
	}
}
