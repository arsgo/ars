package main

import (
	"sort"
	"strings"

	"github.com/colinyl/ars/cluster"
)

//Bind 绑定服务
func (rc *RCServer) Bind() (err error) {
	rc.snap.Path, err = rc.clusterClient.CreateRCServer(rc.snap.GetSnap())
	if err != nil {
		return
	}
	rc.clusterClient.WatchRCServerChange(func(items []*cluster.RCServerItem, err error) {
		rc.IsMaster = rc.IsMasterServer(items, rc.snap.Path)
		if rc.IsMaster {
			rc.snap.Server = SERVER_MASTER
			rc.Log.Info("current server is ", SERVER_MASTER)
		} else {
			rc.Log.Info("current server is ", SERVER_SLAVE)
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
