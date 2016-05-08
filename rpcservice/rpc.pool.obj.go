package rpcservice

import (
	"log"
	"sync"
	"time"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/pool"
)

type rpcServerService struct {
	Status bool
	IP     string
}
type rpcClientFactory struct {
	ip  string
	Log *logger.Logger
}
type RPCServerPool struct {
	pool    *pool.ObjectPool
	servers map[string]*rpcServerService
	lk      sync.Mutex
	Log     *logger.Logger
}

func newRPCClientFactory(ip string, log *logger.Logger) *rpcClientFactory {
	return &rpcClientFactory{ip: ip, Log: log}
}
func (j *rpcClientFactory) create() (p pool.Object, err error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	o := NewRPCClient(j.ip)
	err = o.Open()
	p = o
	return
}
func (j *rpcClientFactory) Create() (p pool.Object, err error) {
	for {
		p, err = j.create()
		if err == nil {
			return
		}
		time.Sleep(time.Second * 5)
	}
	return
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
