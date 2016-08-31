package rpcservice

import (
	"errors"
	"time"

	"github.com/arsgo/lib4go/logger"
	"github.com/arsgo/lib4go/pool"
)

type rpcClientFactory struct {
	ip         string
	loggerName string
	closeQueue chan int
	Log        logger.ILogger
	isClose    bool
}

func newRPCClientFactory(ip string, loggerName string) *rpcClientFactory {
	rf := &rpcClientFactory{ip: ip, loggerName: loggerName, closeQueue: make(chan int, 1)}
	rf.Log, _ = logger.New(loggerName)
	return rf
}
func (fac *rpcClientFactory) Create() (p pool.Object, err error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	if fac.isClose {
		err = errors.New("factory is closed")
		return
	}
	client := NewRPCClientTimeout(fac.ip, time.Second*5, fac.loggerName)
	err = client.Open()
	if err != nil {
		return
	}
	p = client
	return
}
func (fac *rpcClientFactory) Close() {
	fac.isClose = true
	fac.closeQueue <- 1
}
