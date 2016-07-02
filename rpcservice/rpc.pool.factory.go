package rpcservice

import (
	"errors"

	"github.com/colinyl/lib4go/pool"
)

type rpcClientFactory struct {
	ip         string
	loggerName string
	closeQueue chan int
	isClose    bool
}

func newRPCClientFactory(ip string, loggerName string) *rpcClientFactory {
	return &rpcClientFactory{ip: ip, loggerName: loggerName, closeQueue: make(chan int, 1)}
}
func (j *rpcClientFactory) Create() (p pool.Object, err error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	if j.isClose {
		err = errors.New("factory is closed")
		return
	}
	ch := make(chan *RpcClientConn, 1)
	Subscribe(j.ip, ch, j.loggerName)
	select {
	case client := <-ch:
		p = client.Client
		err = client.Err
	}
	return

}
func (j *rpcClientFactory) Close() {
	j.isClose = true
	j.closeQueue <- 1
}
