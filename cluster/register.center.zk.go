package cluster

import (
	"sort"
	"strings"
)

func (d *rcServer) createRCServer(path string, value string) (err error) {
	d.Path, err = d.zkClient.ZkCli.CreateSeqNode(path, value)
	if err != nil {
		return
	}
	d.setOnlineParams(d.isMaster())
	return
}

func (d *rcServer) resetRCSnap() (err error) {
	err = d.zkClient.ZkCli.UpdateValue(d.Path, d.snap.GetSnap())
	return
}

func (d *rcServer) isMaster() bool {
	servers, _ := d.zkClient.ZkCli.GetChildren(d.rcServerRoot)
	sort.Sort(sort.StringSlice(servers))
	return len(servers) == 0 || strings.HasSuffix(d.Path, servers[0])
}

func (d *rcServer) setOnlineParams(master bool) {
	d.IsMasterServer = master
	d.snap.Path = d.Path
	if d.IsMasterServer {
		d.snap.Server = SERVER_MASTER
		d.Server = SERVER_MASTER
	} else {
		d.snap.Server = SERVER_SLAVE
		d.Server = SERVER_SLAVE
	}
	d.Log.Infof("current server is %s", d.Server)
}
