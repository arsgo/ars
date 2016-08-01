package main

import (
	"fmt"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/colinyl/ars/cluster"
)

//BindRCServer 绑定服务
func (rc *RCServer) BindRCServer() (err error) {
	rc.Log.Info("------bind rcserver")
	rc.snap.Address = fmt.Sprint(rc.snap.ip, rc.rcRPCServer.Address)
	rc.snap.Path, err = rc.clusterClient.CreateRCServer(rc.snap.GetSnap())
	if err != nil {
		return
	}
	rc.clusterClient.UpdateSnap(rc.snap.Path, rc.snap.GetSnap())
	rc.clusterClient.WatchRCServerChange(func(items []*cluster.RCServerItem, err error) {
		isMaster := rc.IsMasterServer(items)
		if isMaster && !rc.IsMaster {
			rc.IsMaster = true
			rc.snap.Server = SERVER_MASTER
			rc.Log.Info("::current server is ", rc.snap.Server)

			go rc.clusterClient.WatchJobConfigChange(func(config map[string]cluster.JobItem, err error) {
				rc.BindJobScheduler(config, err)
			})
			go rc.clusterClient.WatchServiceProviderChange(func(lst cluster.ServiceProviderList, err error) {
				//重新发布服务
				rc.Log.Info(" |-> rpc service provider changed")
				rc.currentServices.Set("*", lst)
				rc.PublishNow()
				rc.startSync.Done("INIT.SERVER")
			})
			go rc.clusterClient.WatchRCTaskChange(func(task cluster.RCServerTask, err error) {
				if err != nil {
					rc.Log.Error(err)
					return
				}
				rc.spRPCClient.SetPoolSize(task.RPCPoolSetting.MinSize, task.RPCPoolSetting.MaxSize)
				rc.BindCrossAccess(task)
			})
		} else if !isMaster {
			rc.IsMaster = false
			rc.snap.Server = SERVER_SLAVE
			rc.Log.Info("::current server is ", rc.snap.Server)
			go rc.clusterClient.WatchRCTaskChange(func(task cluster.RCServerTask, err error) {
				rc.spRPCClient.SetPoolSize(task.RPCPoolSetting.MinSize, task.RPCPoolSetting.MaxSize)
				rc.startSync.Done("INIT.SERVER")
			})
		}
	})
	rc.startSync.WaitAndAdd(1)
	rc.clusterClient.WatchRPCServiceChange(func(services map[string][]string, err error) {
		defer rc.startSync.Done("INIT.SRV.CNG")
		rc.BindSPServers(services, err)
	})
	return
}

//PublishAll 发布所有服务
func (rc *RCServer) PublishAll() {
	currentServices, err := rc.clusterClient.GetServiceProviders()
	if err != nil {
		rc.Log.Error(err)
		return
	}
	rc.currentServices.Set("*", currentServices)
	crossClusters := rc.crossDomain.GetAll()
	for domain, clt := range crossClusters {
		client := clt.(cluster.IClusterClient)
		crossService, err := client.GetServiceProviders()
		if err != nil {
			rc.Log.Error(err)
			continue
		}
		rc.crossService.Set(domain, crossService)
	}
	rc.PublishNow()
}

//BindSPServers 绑定service provider servers
func (rc *RCServer) BindSPServers(services map[string][]string, err error) {
	if err != nil {
		return
	}
	rc.Log.Info(" |-> rpc services changed")
	ip := rc.spRPCClient.ResetRPCServer(services)
	tasks, er := rc.clusterClient.FilterRPCService(services)
	if er != nil {
		rc.Log.Error(er)
		return
	}
	if rc.rcRPCServer.UpdateTasks(tasks) > 0 {
		rc.Log.Info(" |-> local services has changed:", len(tasks), ",", ip)
	}
}

//PublishNow 立即发布服务
func (rc *RCServer) PublishNow() {
	defer func() {
		if r := recover(); r != nil {
			rc.Log.Fatal(r, string(debug.Stack()))
		}
	}()
	//立即发布服务
	services := rc.MergeService()
	rc.Log.Infof("->publish services:%d", len(services))
	err := rc.clusterClient.PublishRPCServices(services)
	if err != nil {
		rc.Log.Error(err)
	}
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
func (rc *RCServer) MergeService() (lst cluster.ServiceProviderList) {
	lst = make(map[string][]string)
	services := rc.currentServices.Get("*")
	if services != nil {
		lst = services.(cluster.ServiceProviderList)
	}
	crossServices := rc.crossService.GetAll()
	for domain, svs := range crossServices {
		service := svs.(cluster.ServiceProviderList)
		for i, v := range service {
			if len(v) > 0 {
				lst[i+"@"+domain] = v
			}
		}
	}
	rc.Log.Info("merge.crossServices:", crossServices)
	rc.Log.Info("merge.all:", lst)
	return lst
}
