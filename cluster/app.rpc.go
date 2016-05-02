package cluster

func (a *appServer) BindRCServer(configs []*RCServerConfig, err error) error {
	servers := make(map[string]string)
	services := make(map[string][]string)
	services["-"] = make([]string, 0)
	for _, v := range configs {
		sv := v.Address
		servers[sv] = sv
		services["-"] = append(services["-"], sv)
	}
	a.rcServicesMap.setData(services)
	a.rcServerPool.Register(servers)
	return nil
}
