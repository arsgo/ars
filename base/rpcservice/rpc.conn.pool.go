//用于管理连接池信息，所有连接都通过同一个work线程创建，当创建成功后通过chan返回连接

package rpcservice

import (
	"fmt"
	"runtime/debug"
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
	notify     chan *RpcClientConn
	loggerName string
}
type worker struct {
	address     string
	status      chan bool
	connect     int
	subscribers chan *subscriber
	loggerName  string
	Log         logger.ILogger
}

type connPool struct {
	workers    map[string]*worker
	mutex      sync.Mutex
	status     bool
	loggerName string
}

var globalPool *connPool

func init() {
	globalPool = NewConnPool()
}

func Subscribe(address string, notify chan *RpcClientConn, loggerName string) {
	globalPool.Subscribe(address, notify, loggerName)
}

func NewConnPool() (conn *connPool) {
	conn = &connPool{workers: make(map[string]*worker)}

	return
}
func (n *connPool) recover() {
	if r := recover(); r != nil {
		log, _ := logger.Get("sys/conn.pool")
		log.Fatal(r, string(debug.Stack()))
	}
}
func (n *connPool) Subscribe(address string, notify chan *RpcClientConn, loggerName string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	wkr, ok := n.workers[address]
	if !ok {
		wkr = &worker{address: address, status: make(chan bool, 1), loggerName: loggerName}
		wkr.Log, _ = logger.Get(loggerName)
		wkr.subscribers = make(chan *subscriber, 100)
		wkr.status <- true
		n.workers[address] = wkr
		go func() {
			n.recover()
			wkr.doWork()
		}()
	}
	wkr.subscribers <- &subscriber{notify: notify, loggerName: loggerName}
}
func (w *worker) doWork() {
	tp := time.NewTicker(time.Second * 3)
	for {
		select {
		case sub := <-w.subscribers:
			{
				if w.connect == c_cant_connect {
					sub.notify <- &RpcClientConn{Err: fmt.Errorf("cant connect server:%s", w.address)}
				} else {
					//w.Log.Info(" -> connect to:", w.address)
					client := NewRPCClientTimeout(w.address, time.Second*5, sub.loggerName)
					err := client.Open()
					if err == nil {
						w.connect = c_connecdted
					} else {
						w.Log.Error(err)
						w.connect = c_cant_connect
					}
					sub.notify <- &RpcClientConn{Client: client, Err: err}
				}
			}
		case <-tp.C:
			if w.connect == c_cant_connect {
				w.Log.Info(" -> 定时重连:", w.address)
				client := NewRPCClient(w.address, w.loggerName)
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
