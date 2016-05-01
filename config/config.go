package config

type sysConfig struct {
	ZKServers []string
	Domain    string
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
	c.ZKServers = []string{"171.221.206.81:2181"}
	c.Domain = "/grs/core"
	return c
}
