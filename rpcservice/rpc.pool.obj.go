package rpcservice

import (
	"fmt"
	"log"
	"sync"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/pool"
)

type rpcServerService struct {
	Status bool
	IP     string
}
type rpcClientFactory struct {
	ip string
}
type RPCServerPool struct {
	pool    *pool.ObjectPool
	servers map[string]*rpcServerService
	lk      sync.Mutex
	Log     *logger.Logger
}

func newRPCClientFactory(ip string) *rpcClientFactory {
	log.Println(ip)
	return &rpcClientFactory{ip: ip}
}

func (j *rpcClientFactory) Create() (pool.Object, error) {
	o := NewRPCClient(j.ip)
	if err := o.Open(); err != nil {
	
		fmt.Println(err.Error())
		return nil, err
	}
	return o, nil
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
