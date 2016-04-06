package rpcservice

import (
	"net"
	"strings"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/colinyl/ars/rpcservice/rpc"
)

type rpcClient struct {
	Address   string
	transport thrift.TTransport
	client    *rpc.ServiceProviderClient
	isFatal   bool
}


func NewRPCClient(address string) *rpcClient {
	addr := address
	if !strings.Contains(address, ":") {
		addr = net.JoinHostPort(address, "1016")
	}
	return &rpcClient{Address: addr}
}

func (client *rpcClient) Open() (err error) {
	client.transport, err = thrift.NewTSocket(client.Address)
	if err != nil {
		return err
	}

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	pf := thrift.NewTBinaryProtocolFactoryDefault()

	useTransport := transportFactory.GetTransport(client.transport)
	client.client = rpc.NewServiceProviderClientFactory(useTransport, pf)
	if err := client.client.Transport.Open(); err != nil {
		return err
	}
	return nil
}

func (j *rpcClient) Request(name string, input string) (string, error) {
	return j.client.Request(name, input)
}

func (j *rpcClient) Send(name string, input string, data []byte) (string, error) {
	return j.client.Send(name, input, data)
}
func (j *rpcClient) Get(name string, input string) ([]byte, error) {
	return j.client.Get(name, input)
}
func (j *rpcClient) Close() {
	j.transport.Close()
}

func (j *rpcClient) Check() bool {
	return !j.isFatal && j.transport != nil
}
func (j *rpcClient) Fatal() {
	j.isFatal = true
}
