package cluster

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"github.com/colinyl/lib4go/utility"
)

func checkIP(origin string) bool {
	ips := fmt.Sprintf(",%s,", origin)
	llocal := fmt.Sprintf(",%s,", zkClient.LocalIP)
	return strings.Contains(ips, llocal)
}

func getSPConfig(path string) (svs []*spService, err error) {
	values, err := zkClient.ZkCli.GetValue(path)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(values), &svs)
	return
}

func (d *spServer) getNewDataMap(serviceName string) *utility.DataMap {
	nmap := d.dataMap.Copy()
	nmap.Set("serviceName", serviceName)
	nmap.Set("last", fmt.Sprintf("%d", time.Now().Unix()))
	return nmap
}
func (d *spServer) createSPPath(path string, nmap *utility.DataMap) {
	if zkClient.ZkCli.Exists(path) {
		return
	}
	_, err := zkClient.ZkCli.CreateTmpNode(path, nmap.Translate(serviceProviderValue))
	if err != nil {
		return
	}
}
func (d *spServer) deleteSPPath(path string) {
	if zkClient.ZkCli.Exists(path) {
		zkClient.ZkCli.Delete(path)
	}
}
func (d *spServer) groupService() (aloneService map[string]*spService,
	sharedService map[string]*spService) {
	aloneService = make(map[string]*spService)
	sharedService = make(map[string]*spService)
	svs, _ := getSPConfig(d.serviceConfig)
	for _, v := range svs {
		if strings.EqualFold(v.Mode, eModeAlone) {
			aloneService[v.Name] = v
		} else if strings.EqualFold(v.Mode, eModeShared) {
			sharedService[v.Name] = v
		}
	}
	return
}
