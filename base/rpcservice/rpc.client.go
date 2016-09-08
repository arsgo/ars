package rpcservice

import (
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/base/rpcservice/rpc"
	"github.com/arsgo/lib4go/logger"
)

type RPCClient struct {
	Address     string
	socket      *thrift.TSocket
	transport   thrift.TTransport
	client      *rpc.ServiceProviderClient
	sync        base.Sync
	lastUseTime time.Time
	closeChan   chan int
	isAvailable bool
	isClose     bool
	connTimeout time.Duration
	rwTimeout   time.Duration
	Log         logger.ILogger
	lk          sync.Mutex
}

func NewRPCClient(address string, loggerName string) (client *RPCClient) {
	return NewRPCClientTimeout(address, time.Second*5, loggerName)
}
func NewRPCClientTimeout(address string, timeout time.Duration, loggerName string) (client *RPCClient) {
	addr := address
	if !strings.Contains(address, ":") {
		addr = net.JoinHostPort(address, "1016")
	}
	client = &RPCClient{Address: addr, connTimeout: timeout, rwTimeout: timeout, lastUseTime: time.Now(), isClose: false}
	client.sync = base.NewSync(1)
	client.closeChan = make(chan int, 1)
	client.Log, _ = logger.Get(loggerName)
	//go client.sendHeartbeat()
	return
}
func (n *RPCClient) Available() bool {
	return n.isAvailable
}
func (n *RPCClient) recover() {
	if r := recover(); r != nil {
		n.Log.Fatal(r, string(debug.Stack()))
	}
}

func (client *RPCClient) Open() (err error) {
	return client.OpenTimeout(client.connTimeout)
}
func (client *RPCClient) SetRWTimeout(timeout time.Duration) {
	client.rwTimeout = timeout
	client.socket.SetTimeout(client.connTimeout, client.rwTimeout, client.rwTimeout)
}
func (client *RPCClient) OpenTimeout(timeout time.Duration) (err error) {
	defer client.recover()
	client.lastUseTime = time.Now()
	client.socket, err = thrift.NewTSocketTimeout(client.Address, timeout)
	if err != nil {
		return errors.New(fmt.Sprint("new client error:", client.Address, ",", err.Error()))
	}
	client.transport = client.socket
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	pf := thrift.NewTBinaryProtocolFactoryDefault()

	useTransport := transportFactory.GetTransport(client.transport)
	client.client = rpc.NewServiceProviderClientFactory(useTransport, pf)
	if err := client.client.Transport.Open(); err != nil {
		return errors.New(fmt.Sprint("open client error :", client.Address, ",", err.Error()))
	}
	client.isAvailable = true
	client.sync.Done("FIRST_OPEN")
	return nil
}

func (client *RPCClient) Request(name string, input string, session string, timeout time.Duration) (r string, e error) {
	//client.lk.Lock()
	//defer client.lk.Unlock()
	defer client.recover()
	client.lastUseTime = time.Now()
	r, e = client.client.Request(name, input, session, timeout.Nanoseconds())
	if e != nil {
		client.isAvailable = false
		return
	}
	client.lastUseTime = time.Now()
	return
}

func (client *RPCClient) Send(name string, input string, data []byte, timeout time.Duration) (string, error) {
	client.lastUseTime = time.Now()
	return client.client.Send(name, input, data, timeout.Nanoseconds())
}
func (client *RPCClient) Get(name string, input string, timeout time.Duration) ([]byte, error) {
	client.lastUseTime = time.Now()
	return client.client.Get(name, input, timeout.Nanoseconds())
}

func (client *RPCClient) Heartbeat(input string, timeout time.Duration) (r string, err error) {
	client.lastUseTime = time.Now()
	return client.client.Heartbeat(input, timeout.Nanoseconds())
}

func (client *RPCClient) Close() {
	defer client.recover()
	client.isClose = true
	select {
	case client.closeChan <- 1:
	default:
	}
	if client.transport != nil {
		client.transport.Close()
	}

}
func (client *RPCClient) sendHeartbeat() {
	client.sync.Wait()
	tk := time.NewTicker(time.Second * 3)
START:
	for {
		select {
		case <-client.closeChan:
			break START
		case <-tk.C:
			if client.isClose {
				break START
			}
			client.lk.Lock()
			if time.Now().Sub(client.lastUseTime).Seconds() > 3 {
				r, err := client.Heartbeat("hb", time.Second)
				client.isAvailable = err == nil && strings.EqualFold(r, "success")
				if !client.isAvailable {
					client.Close()
					client.OpenTimeout(time.Second)
				}
			}
			client.lk.Unlock()
		}
	}
}
