package main

import (
	"sort"
	"strings"

	"github.com/colinyl/ars/cluster"
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
