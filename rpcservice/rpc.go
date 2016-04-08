package rpcservice

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/colinyl/lib4go/logger"
)

type rpcHandler interface {
	Request(name string, input string) (r string, err error)
	Send(name string, input string, data []byte) (r string, err error)
	Get(name string, input string) (data []byte, err error)
}

//JobProviderServer
type RPCServer struct {
	Address string
	Handler rpcHandler
	log     *logger.Logger
	server  *thrift.TSimpleServer
}
