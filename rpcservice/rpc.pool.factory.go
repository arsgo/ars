package rpcservice

import (
	"fmt"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/pool"
)

type rpcClientFactory struct {
	ip         string
	Log        *logger.Logger
	closeQueue chan int
	isClose    bool
}

func newRPCClientFactory(ip string, log *logger.Logger) *rpcClientFactory {
	return &rpcClientFactory{ip: ip, Log: log, closeQueue: make(chan int, 1)}
}
func (j *rpcClientFactory) Create() (p pool.Object, err error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	ch := make(chan *RpcClientConn, 1)
	fmt.Println("add subscribe")
	Subscribe(j.ip, ch)
	fmt.Println("wait connect create")
	for {
		select {
		case client := <-ch:
			fmt.Println("recv subscriber response")
			p = client.Client
			err = client.Err
			return
		case <-j.closeQueue:
			return nil, nil
		}
	}
}
func (j *rpcClientFactory) Close() {
	j.isClose = true
	j.closeQueue <- 1
}
