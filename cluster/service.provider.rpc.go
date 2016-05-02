package cluster

import (
	"fmt"

	"github.com/colinyl/ars/rpcservice"
)

func (d *spServer) StartRPC() error {
	port := rpcservice.GetLocalRandomAddress()
	d.dataMap.Set("port",port)
	d.snap.Address = fmt.Sprintf("%s%s", d.zkClient.LocalIP, port)
	d.rpcServer = rpcservice.NewRPCServer(port, NewScript(d))
	return d.rpcServer.Serve()
}
