package config

import "github.com/colinyl/lib4go/utility"

type sysConfig struct {
	ZKServers []string
	Domain    string
	Mask      []string
	IP        string
}

var _config *sysConfig

func init() {
	_config = getDefConfig()
}

func Get() *sysConfig {
	return _config
}

func readConfig() *sysConfig {
	return &sysConfig{}
}
func getDefConfig() *sysConfig {
	c := &sysConfig{}
	c.ZKServers = []string{"192.168.101.161:2181"} //"171.221.206.81:2181",
	c.Domain = "/grs/core"
	c.Mask = []string{"192.168", "172.16"}
	c.IP = utility.GetLocalIPAddress(c.Mask...)
	return c
}
