package cluster

import "fmt"

func (a *appServer) BindRCServer(configs []*RCServerConfig, err error) error {
	servers := make(map[string]string)
	services := make(map[string][]string)
	services["-"] = make([]string, 0)
	for _, v := range configs {
		sv := fmt.Sprintf("%s%s", v.IP, v.Port)
		servers[sv] = sv
		services["-"] = append(services["-"], sv)
	}
	a.rcServicesMap.setData(services)
	a.rcServerPool.Register(servers)
    a.Log.Infof("app.script.send.group:%s\r\n",a.rcServicesMap.Next("-"))
    
	return nil
}
