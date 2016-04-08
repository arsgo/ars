package cluster

import "github.com/colinyl/ars/rpcservice"

func (d *spServer) StartRPC() {
	address := rpcservice.GetLocalRandomAddress()
	d.Port = address
	d.dataMap.Set("port", d.Port)
	d.rpcServer = rpcservice.NewRPCServer(address, NewScript(d))
	d.rpcServer.Serve()
}
