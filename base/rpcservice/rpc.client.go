package rpcservice

import (
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/base/rpcservice/rpc"
	"github.com/arsgo/lib4go/logger"
)

type RPCClient struct {
	Address     string
	transport   thrift.TTransport
	client      *rpc.ServiceProviderClient
	lastUseTime time.Time
	closeChan   chan int
	isFatal     bool
	isClose     bool
	timeout     time.Duration
	Log         logger.ILogger
}

func NewRPCClient(address string, loggerName string) (client *RPCClient) {
	return NewRPCClientTimeout(address, time.Second*30, loggerName)
}
func NewRPCClientTimeout(address string, timeout time.Duration, loggerName string) (client *RPCClient) {
	addr := address
	if !strings.Contains(address, ":") {
		addr = net.JoinHostPort(address, "1016")
	}
	client = &RPCClient{Address: addr, timeout: time.Second * 30, lastUseTime: time.Now()}
	client.closeChan = make(chan int, 1)
	client.Log, _ = logger.Get(loggerName)
	return
}
func (n *RPCClient) recover() {
	if r := recover(); r != nil {
		n.Log.Fatal(r, string(debug.Stack()))
	}
}
func (client *RPCClient) Open() (err error) {
	defer client.recover()
	client.lastUseTime = time.Now()
	client.transport, err = thrift.NewTSocketTimeout(client.Address, client.timeout)
	if err != nil {
		return errors.New(fmt.Sprint("new client error:", client.Address, ",", err.Error()))
	}

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	pf := thrift.NewTBinaryProtocolFactoryDefault()

	useTransport := transportFactory.GetTransport(client.transport)
	client.client = rpc.NewServiceProviderClientFactory(useTransport, pf)
	if err := client.client.Transport.Open(); err != nil {
		return errors.New(fmt.Sprint("open client error :", client.Address, ",", err.Error()))
	}
//	go client.sendHeartbeat()
	return nil
}

func (client *RPCClient) Request(name string, input string, session string) (r string, e error) {
	defer client.recover()
	defer base.RunTime("rpc request once", time.Now())
	client.lastUseTime = time.Now()
	client.isClose = true
	r, er := client.client.Request(name, input, session)
	if er != nil {
		r = er.Error()
	}
	return
}

func (client *RPCClient) Send(name string, input string, data []byte) (string, error) {
	client.lastUseTime = time.Now()
	return client.client.Send(name, input, data)
}
func (client *RPCClient) Get(name string, input string) ([]byte, error) {
	client.lastUseTime = time.Now()
	return client.client.Get(name, input)
}

func (client *RPCClient) Heartbeat(input string) (r string, err error) {
	client.lastUseTime = time.Now()
	return client.client.Heartbeat(input)
}

func (client *RPCClient) Close() {
	defer client.recover()
	client.closeChan <- 1
	if client.transport != nil {
		client.transport.Close()
	}
}
func (client *RPCClient) sendHeartbeat() {
	tk := time.NewTicker(time.Second * 2)
START:
	for {
		select {
		case <-client.closeChan:
			break START
		case <-tk.C:
			if client.isClose {
				break START
			}
			if time.Now().Sub(client.lastUseTime).Seconds() > 5 {
				r, err := client.Heartbeat("hb")
				if err != nil || !strings.EqualFold(r, "success") {
					client.isFatal = true
				} else {
					client.isFatal = false
				}
			}
		}
	}
}
