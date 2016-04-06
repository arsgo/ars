package cluster

import (
	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/logger"
)

type rcServerRPCHandler struct {
	server *rcServer
	Log    *logger.Logger
}

func (r *rcServerRPCHandler) Request(name string, input string) (string, error) {
	r.Log.Infof("recv request:%s", name)
	sv := r.server.spServicesMap.get(name)
	return r.server.spServerPool.Request(sv, input)

}
func (r *rcServerRPCHandler) Send(name string, input string, data []byte) (string, error) {
	return "", nil
}
func (r *rcServerRPCHandler) Get(name string, input string) ([]byte, error) {
	return make([]byte, 0), nil
}

func (d *rcServer) StartRPCServer() {
	address := rpcservice.GetLocalRandomAddress()
	d.Port = address
	d.dataMap.Set("port", d.Port)
	d.rpcServer = rpcservice.NewRPCServer(address, &rcServerRPCHandler{server: d, Log: d.Log})
	d.rpcServer.Serve()
	d.resetRCSnap()
}

func (r *rcServer) BindSPServer(services map[string][]string) {
	r.spServicesMap.setData(services)
	r.spServerPool.Register(getServers(services))
}

func getServers(services map[string][]string) map[string]string {
	servers := make(map[string]string)
	for _, v := range services {
		for _, value := range v {
			if _, ok := servers[value]; !ok {
				servers[value] = value
			}
		}
	}
	return servers
}
