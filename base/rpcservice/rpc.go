package rpcservice

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/arsgo/lib4go/logger"
)

type rpcHandler interface {
	Request(name string, input string, session string) (r string, err error)
	Send(name string, input string, data []byte) (r string, err error)
	Get(name string, input string) (data []byte, err error)
	Heartbeat(input string) (r string, err error)
}

type RPCServer struct {
	Address string
	Handler rpcHandler
	log     logger.ILogger
	server  *thrift.TSimpleServer
}
