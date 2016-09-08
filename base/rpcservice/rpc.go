package rpcservice

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/arsgo/lib4go/logger"
)

type rpcHandler interface {
	Request(name string, input string, session string,timeout int64) (r string, err error)
	Send(name string, input string, data []byte,timeout int64) (r string, err error)
	Get(name string, input string,timeout int64) (data []byte, err error)
	Heartbeat(input string,timeout int64) (r string, err error)
}

type RPCServer struct {
	Address string
	Handler rpcHandler
	log     logger.ILogger
	server  *thrift.TSimpleServer
}
