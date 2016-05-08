package rpcservice

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/colinyl/ars/rpcservice/rpc"
)

type RPCClient struct {
	Address   string
	transport thrift.TTransport
	client    *rpc.ServiceProviderClient
	isFatal   bool
}

func NewRPCClient(address string) *RPCClient {
	addr := address
	if !strings.Contains(address, ":") {
		addr = net.JoinHostPort(address, "1016")
	}
	return &RPCClient{Address: addr}
}

func (client *RPCClient) Open() (err error) {
	client.transport, err = thrift.NewTSocket(client.Address)
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

func (j *RPCClient) Request(name string, input string) (string, error) {
	return j.client.Request(name, input)
}

func (j *RPCClient) Send(name string, input string, data []byte) (string, error) {
	return j.client.Send(name, input, data)
}
func (j *RPCClient) Get(name string, input string) ([]byte, error) {
	return j.client.Get(name, input)
}
func (j *RPCClient) Close() {
	defer func() {
		recover()
	}()
	if j.transport != nil {
		j.transport.Close()
	}

}

func (j *RPCClient) Check() bool {
	return !j.isFatal && j.transport != nil
}
func (j *RPCClient) Fatal() {
	j.isFatal = true
}
