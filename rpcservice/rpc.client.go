package rpcservice

import (
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
		fmt.Printf("new client error:%s,%s\r\n", client.Address, err.Error())
		return err
	}

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	pf := thrift.NewTBinaryProtocolFactoryDefault()

	useTransport := transportFactory.GetTransport(client.transport)
	client.client = rpc.NewServiceProviderClientFactory(useTransport, pf)
	if err := client.client.Transport.Open(); err != nil {
		fmt.Printf("open client error :%s,%s",client.Address,err.Error())
		return err
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
