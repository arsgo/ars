package config

type sysConfig struct {
	ZKServers       []string
	Domain   string
}
var _config *sysConfig
func init(){
    _config=getDefConfig()
}


func Get() *sysConfig{
    return _config
}

func readConfig() *sysConfig {
	return &sysConfig{}
}
func getDefConfig() *sysConfig {
	c := &sysConfig{}
	c.ZKServers = []string{"192.168.101.161:2181"}
	c.Domain = "/grs/pay"
	return c
}
