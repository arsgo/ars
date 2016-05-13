package rpcservice

import (
	"log"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/pool"
)

type rpcServerService struct {
	Status bool
	IP     string
}

type RPCServerPool struct {
	pool    *pool.ObjectPool
	servers map[string]*rpcServerService
	Log     *logger.Logger
}

func NewRPCServerPool() *RPCServerPool {
	var err error
	pl := &RPCServerPool{}
	pl.pool = pool.New()
	go pl.autoClearUp()
	pl.servers = make(map[string]*rpcServerService)
	pl.Log, err = logger.New("rc server", true)
	if err != nil {
		log.Println(err)
	}
	return pl
}
