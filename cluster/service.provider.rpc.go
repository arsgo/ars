package cluster

import (
	"fmt"

	"github.com/colinyl/ars/rpcservice"
)

func (d *spServer) StartRPC() error {
	port := rpcservice.GetLocalRandomAddress()
	d.dataMap.Set("port", port)
	d.snap.Address = fmt.Sprintf("%s%s", d.zkClient.IP, port)
	d.rpcServer = rpcservice.NewRPCServer(port, d.scriptEngine)
	return d.rpcServer.Serve()
}
