//用于管理连接池信息，所有连接都通过同一个work线程创建，当创建成功后通过chan返回连接

package rpcservice

import (
	"fmt"
	"sync"
	"time"

	"github.com/colinyl/lib4go/logger"
)

const (
	c_not_connect  = 0
	c_connecdted   = 1
	c_cant_connect = 2
)

type RpcClientConn struct {
	Client *RPCClient
	Err    error
}

type subscriber struct {
	notify chan *RpcClientConn
}
type worker struct {
	address     string
	status      chan bool
	connect     int
	subscribers chan *subscriber
	Log         *logger.Logger
}

type connPool struct {
	workers map[string]*worker
	mutex   sync.Mutex
	status  bool
	Log     *logger.Logger
}

var globalPool *connPool

func init() {
	globalPool = NewConnPool()
}

func Subscribe(address string, notify chan *RpcClientConn) {
	globalPool.Subscribe(address, notify)
}

func NewConnPool() (conn *connPool) {
	conn = &connPool{workers: make(map[string]*worker)}
	conn.Log, _ = logger.New("conn pool", true)
	return
}

func (n *connPool) Subscribe(address string, notify chan *RpcClientConn) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	wkr, ok := n.workers[address]
	if !ok {
		wkr = &worker{address: address, status: make(chan bool, 1)}
		wkr.Log = n.Log
		wkr.subscribers = make(chan *subscriber, 100)
		wkr.status <- true
		n.workers[address] = wkr
		go wkr.doWork()
	}
	wkr.subscribers <- &subscriber{notify: notify}
}
func (w *worker) doWork() {
	tp := time.NewTicker(time.Second * 5)
	for {
		select {
		case sub := <-w.subscribers:
			{
				if w.connect == c_cant_connect {
					sub.notify <- &RpcClientConn{Err: fmt.Errorf("cant connect server:%s", w.address)}
				} else {
					//	w.Log.Info(" -> connect to:", w.address)
					client := NewRPCClient(w.address)
					err := client.Open()
					if err == nil {
						w.connect = c_connecdted
					} else {
						w.connect = c_cant_connect
					}
					sub.notify <- &RpcClientConn{Client: client, Err: err}
				}
			}
		case <-tp.C:
			if w.connect == c_cant_connect {
				w.Log.Info(" -> 定时重连:", w.address)
				client := NewRPCClient(w.address)
				err := client.Open()
				if err == nil {
					w.connect = c_connecdted
				} else {
					w.Log.Error(err)
				}
				if w.connect == c_connecdted {
					client.Close()
				}
			}
		}
	}
}
