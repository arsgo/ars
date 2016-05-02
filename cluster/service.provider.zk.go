package cluster

import (
	"strings"

	"github.com/colinyl/lib4go/utility"
)

func (d *spServer) getNewDataMap(serviceName string) utility.DataMap {
	nmap := d.dataMap.Copy()
	nmap.Set("serviceName", serviceName)
	return nmap
}

func (d *spServer) createSPPath(path string, value string) {
	if d.zkClient.ZkCli.Exists(path) {
		return
	}
	_, err := d.zkClient.ZkCli.CreateTmpNode(path, value)
	if err != nil {
		return
	}
}
func (d *spServer) deleteSPPath(path string) {
	if d.zkClient.ZkCli.Exists(path) {
		d.zkClient.ZkCli.Delete(path)
	}
}
func (d *spServer) groupService() (aloneService map[string]spService,
	sharedService map[string]spService) {
	aloneService = make(map[string]spService)
	sharedService = make(map[string]spService)
	svs, _ := d.zkClient.getSPConfig(d.serviceConfig)
	for _, v := range svs {
		if strings.EqualFold(v.Mode, eModeAlone) {
			aloneService[v.Name] = v
		} else if strings.EqualFold(v.Mode, eModeShared) {
			sharedService[v.Name] = v
		}
	}
	return
}
