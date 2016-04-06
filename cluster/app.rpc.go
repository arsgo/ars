package cluster

import "fmt"

func (a *appServer) BindRCServer(configs []*RCServerConfig, err error) error {
	servers := make(map[string]string)
	services := make(map[string][]string)
	services["_"] = make([]string, 0)
	for _, v := range configs {
		sv := fmt.Sprintf("%s%s", v.IP, v.Port)
		servers[sv] = sv
		services["_"] = append(services["_"], sv)
	}
	a.rcServicesMap.setData(services)
	a.rcServerPool.Register(servers)
	return nil
}
