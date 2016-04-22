package cluster

import "github.com/colinyl/ars/rpcservice"

func (d *spServer) StartRPC()(error) {
	address :=":2034"// rpcservice.GetLocalRandomAddress()
	d.Port = address
	d.dataMap.Set("port", d.Port)
	d.rpcServer = rpcservice.NewRPCServer(address, NewScript(d))
	return d.rpcServer.Serve()
}
