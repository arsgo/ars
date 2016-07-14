package rpcservice

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/colinyl/ars/rpcservice/rpc"
	"github.com/colinyl/lib4go/logger"
)

type RPCClient struct {
	Address   string
	transport thrift.TTransport
	client    *rpc.ServiceProviderClient
	isFatal   bool
	timeout   time.Duration
	Log       logger.ILogger
}

func NewRPCClient(address string, loggerName string) (client *RPCClient) {
	return NewRPCClientTimeout(address, time.Second*30, loggerName)
}
func NewRPCClientTimeout(address string, timeout time.Duration, loggerName string) (client *RPCClient) {
	addr := address
	if !strings.Contains(address, ":") {
		addr = net.JoinHostPort(address, "1016")
	}
	client = &RPCClient{Address: addr, timeout: time.Second * 30}
	client.Log, _ = logger.Get(loggerName, true)
	return
}
func (n *RPCClient) recover() {
	if r := recover(); r != nil {
		n.Log.Fatal(r)
	}
}
func (client *RPCClient) Open() (err error) {
	defer client.recover()
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
	return nil
}

func (j *RPCClient) Request(name string, input string, session string) (r string, e error) {
	defer j.recover()
	r, er := j.client.Request(name, input, session)
	if er != nil {
		r = er.Error()
	}
	return
}

func (j *RPCClient) Send(name string, input string, data []byte) (string, error) {
	defer j.recover()
	return j.client.Send(name, input, data)
}
func (j *RPCClient) Get(name string, input string) ([]byte, error) {
	defer j.recover()
	return j.client.Get(name, input)
}
func (j *RPCClient) Close() {
	defer j.recover()
	if j.transport != nil {
		j.transport.Close()
	}

}
