package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/arsgo/ars/cluster"
)

//BindRCServer 绑定服务
func (rc *RCServer) BindRCServer() (err error) {
	rc.snap.Address = fmt.Sprint(rc.snap.ip, rc.rcRPCServer.Address)
	rc.snap.Path, err = rc.clusterClient.CreateRCServer(rc.snap.GetSnap())
	if err != nil {
		return
	}

	rc.clusterClient.WatchRCServerChange(func(items []*cluster.RCServerItem, err error) {
		isMaster := rc.IsMasterServer(items)
		if isMaster && !rc.IsMaster {
			rc.IsMaster = true
			rc.snap.Server = SERVER_MASTER
			rc.Log.Info(" -> 当前服务是 [", rc.snap.Server, "]")
			go rc.clusterClient.WatchRCTaskChange(func(task cluster.RCServerTask, err error) {
				if err != nil {
					return
				}
				rc.spRPCClient.SetPoolSize(task.RPCPoolSetting.MinSize, task.RPCPoolSetting.MaxSize)
				rc.BindCrossAccess(task)
			})
			go rc.clusterClient.WatchJobConfigChange(func(config map[string]cluster.JobItem, err error) {
				rc.BindJobScheduler(config, err)
			})
			go rc.clusterClient.WatchSPServerChange(func(lst cluster.RPCServices, err error) {
				defer rc.startSync.Done("INIT.SERVER")
				if err != nil {
					return
				}
				rc.currentServices.Set("*", lst)
				rc.timerPublishServices.Push("publish services")
			})

		} else if !isMaster {
			rc.IsMaster = false
			rc.snap.Server = SERVER_SLAVE
			rc.Log.Info(" -> 当前服务是 [", rc.snap.Server, "]")
			go rc.clusterClient.WatchRCTaskChange(func(task cluster.RCServerTask, err error) {
				defer rc.startSync.Done("INIT.SERVER")
				rc.spRPCClient.SetPoolSize(task.RPCPoolSetting.MinSize, task.RPCPoolSetting.MaxSize)
			})
		}
	})
	rc.startSync.WaitAndAdd(1)
	rc.clusterClient.WatchRPCServiceChange(func(services map[string][]string, err error) {
		defer rc.startSync.Done("INIT.SRV.CNG")
		rc.BindServices(services, err)
	})

	return
}

//resetLoalService 重置本地所有
func (rc *RCServer) resetLoalService() {
	currentServices, err := rc.clusterClient.GetSPServerServices()
	if err != nil {
		return
	}
	rc.currentServices.Set("*", currentServices)
	crossClusters := rc.crossDomain.GetAll()
	for domain, clt := range crossClusters {
		client := clt.(cluster.IClusterClient)
		crossService, err := client.GetSPServerServices()
		if err != nil {
			continue
		}
		rc.crossServices.Set(domain, crossService)
	}
}

//BindServices 绑定services
func (rc *RCServer) BindServices(services map[string][]string, err error) {
	if err != nil {
		return
	}
	rc.spRPCClient.ResetRPCServer(services)
	tasks, er := rc.clusterClient.GetLocalServices(services)
	if er != nil {
		rc.Log.Info(" -> 获取本地服务出错:", er)
		return
	}
	if c := rc.rcRPCServer.UpdateTasks(tasks); c > 0 {
		rc.Log.Info(" -> 本地服务已更新")
		rc.snapLogger.Infof("--------------------services-----------------\n\t\t\t\t\t  %+v\n\t\t\t\t\t  ----------------------------------------------%+v\n\t\t\t\t\t  ----------------------------------------------",
			rc.rcRPCServer.GetServices(), services)
	}
	//else {
	//rc.Log.Infof(" -> 本地无更新:%v, %v", services, rc.rcRPCServer.GetServices())
	//	}
}

//PublishNow 立即发布服务
func (rc *RCServer) PublishNow(p ...interface{}) {
	defer rc.recover()
	//立即发布服务
	services := rc.MergeService()
	rc.clusterClient.PublishServices(services)
}

//IsMasterServer 检查当前RC Server是否是Master
func (rc *RCServer) IsMasterServer(items []*cluster.RCServerItem) bool {
	var servers []string
	for _, v := range items {
		servers = append(servers, v.Path)
	}
	sort.Sort(sort.StringSlice(servers))
	return len(servers) == 0 || strings.EqualFold(rc.snap.Path, servers[0])
}

//MergeService 合并所有服务
func (rc *RCServer) MergeService() (lst cluster.RPCServices) {
	lst = make(map[string][]string)
	services := rc.currentServices.Get("*")
	if services != nil {
		cservices := services.(cluster.RPCServices)
		currentDomain := rc.clusterClient.GetDomainName()
		for name, values := range cservices {
			if strings.HasSuffix(name, currentDomain) {
				lst[name] = values
			}
		}
	}
	crossServices := rc.crossServices.GetAll()
	for _, svs := range crossServices {
		service := svs.(cluster.RPCServices)
		for i, v := range service {
			if len(v) > 0 {
				lst[i] = v
			} else {
				delete(lst, i)
			}
		}
	}
	return lst
}
