package rpcservice

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/colinyl/ars/base/rpcservice/rpc"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/net"
)

func (n *RPCServer) recover() {
	if r := recover(); r != nil {
		n.log.Fatal(r, string(debug.Stack()))
	}
}
func (r *RPCServer) Serve() (er error) {
	defer r.recover()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, er := thrift.NewTServerSocketTimeout(r.Address, time.Hour*24*31)
	if er != nil {
		r.log.Error(er)
		return
	}

	processor := rpc.NewServiceProviderProcessor(r.Handler)
	r.server = thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	r.log.Infof("::start rpc server %s", r.Address)
	go func(r *RPCServer) {
		defer r.recover()
		er = r.server.Serve()
		if er != nil {
			r.log.Error(er)
		}
	}(r)
	time.Sleep(time.Second * 2)
	return
}
func (r *RPCServer) Stop() {
	defer r.recover()
	if r.server != nil {
		r.server.Stop()
	}
}
func NewRPCServer(address string, handler rpcHandler, loggerName string) *RPCServer {
	var err error
	rpcs := &RPCServer{Address: address, Handler: handler}
	rpcs.log, err = logger.Get(loggerName)
	if err != nil {
		log.Println(err)
	}
	return rpcs
}

func GetLocalRandomAddress(start ...int) string {
	return fmt.Sprintf(":%d", getPort(start...))
}

func getPort(start ...int) int {
	s := 10160
	if len(start) > 0 {
		s = start[0]
	}
	for i := 0; i < 100; i++ {
		port := s + i*8
		if net.IsTCPPortAvailable(port) {
			return port
		}
	}
	return -1
}
