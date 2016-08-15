package rpcservice

import (
	"errors"
	"time"

	"github.com/arsgo/lib4go/pool"
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
		//	i//f r := recover(); r != nil {
		//return
		//	}
	}()
	if j.isClose {
		err = errors.New("factory is closed")
		return
	}
	//p, err = getConn(j.ip, j.loggerName)
	client := NewRPCClientTimeout(j.ip, time.Second*5, j.loggerName)
	err = client.Open()
	p = client
	return

}
func (j *rpcClientFactory) Close() {
	j.isClose = true
	j.closeQueue <- 1
	removeWorkers()
}
