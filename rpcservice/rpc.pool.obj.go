package rpcservice

import (
	"log"

	"github.com/colinyl/lib4go/concurrent"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/pool"
)

type rpcServerService struct {
	Status bool
	IP     string
}

type RPCServerPool struct {
	pool    *pool.ObjectPool
	servers concurrent.ConcurrentMap
	Log     *logger.Logger
	MinSize int
	MaxSize int
}

func NewRPCServerPool() *RPCServerPool {
	var err error
	pl := &RPCServerPool{}
	pl.pool = pool.New()
	go pl.autoClearUp()
	pl.servers = concurrent.NewConcurrentMap()
	pl.Log, err = logger.New("rc server", true)
	if err != nil {
		log.Println(err)
	}
	return pl
}
