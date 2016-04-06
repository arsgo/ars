package rpcservice

import (
	"fmt"
	"log"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/net"
	"github.com/colinyl/ars/rpcservice/rpc"
)

func (r *RPCServer) Serve() (er error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, er := thrift.NewTServerSocketTimeout(r.Address, time.Hour*24*31)
	if er != nil {
		r.log.Error(er)
		return
	}

	processor := rpc.NewServiceProviderProcessor(r.Handler)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	r.log.Infof("::start rpc server %s", r.Address)
	go func() {
		er = server.Serve()
		if er != nil {
			r.log.Error(er)
		}
	}()
	return
}

func NewRPCServer(address string, handler rpcHandler) *RPCServer {
	var err error
	rpcs := &RPCServer{Address: address, Handler: handler}
	rpcs.log, err = logger.New("rpc server", true)
	if err != nil {
		log.Println(err)
	}
	return rpcs
}

func GetLocalRandomAddress() string {
	return fmt.Sprintf(":%d", getPort())
}

func getPort() int {
	for i := 0; i < 100; i++ {
		port := 1016 + i*8
		if net.IsTCPPortAvailable(port) {
			return port
		}
	}
	return -1
}
