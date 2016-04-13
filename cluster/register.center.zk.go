package cluster

import (
	"fmt"
	"sort"
	"strings"
	"time"
)


func (d *rcServer) createRCServer(path string, value string) (err error) {
	d.Path, err = d.zkClient.ZkCli.CreateSeqNode(path, value)
	if err != nil {
		return
	}
	d.setOnlineParams(d.isMaster())
	d.resetRCSnap()
	return
}

func (d *rcServer) resetRCSnap() (err error) {
	err = d.zkClient.ZkCli.UpdateValue(d.Path, d.dataMap.Translate(rcServerValue))
	return
}

func (d *rcServer) isMaster() bool {
	servers, _ := d.zkClient.ZkCli.GetChildren(d.rcServerRoot)
	sort.Sort(sort.StringSlice(servers))
	return len(servers) == 0 || strings.HasSuffix(d.Path, servers[0])
}
func (d *rcServer) setLastParams() {
	d.LastPublish = time.Now().Unix()
	d.dataMap.Set("last", fmt.Sprintf("%d", d.LastPublish))
}

func (d *rcServer) setServiceParams() {
	d.LastPublish = time.Now().Unix()
	d.dataMap.Set("last", fmt.Sprintf("%d", d.LastPublish))
	d.dataMap.Set("pst", fmt.Sprintf("%d", d.LastPublish))
}

func (d *rcServer) setOnlineParams(master bool) {
	d.IsMasterServer = master
	d.OnlineTime = time.Now().Unix()
	d.dataMap.Set("path", d.Path)
	d.dataMap.Set("online", fmt.Sprintf("%d", d.OnlineTime))
	d.dataMap.Set("last", fmt.Sprintf("%d", d.OnlineTime))
	if d.IsMasterServer {
		d.dataMap.Set("type", "master")
		d.Server = "master"
	} else {
		d.dataMap.Set("type", "slave")
		d.Server = "slave"
	}
    d.Log.Infof("current server is %s",d.Server)
}
