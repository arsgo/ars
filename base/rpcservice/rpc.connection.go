//用于管理连接池信息，所有连接都通过同一个work线程创建，当创建成功后通过chan返回连接

package rpcservice

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/arsgo/lib4go/logger"
)

type conn struct {
	client *RPCClient
	err    error
}

type worker struct {
	address    string
	close      chan int
	create     chan int
	conns      chan *conn
	connect    int
	loggerName string
	available  int32
	log        logger.ILogger
}

var workers map[string]*worker
var mutex sync.Mutex

func init() {
	workers = make(map[string]*worker, 200)
}

func getConn(address string, loggerName string) (client *RPCClient, err error) {
	mutex.Lock()
	if _, ok := workers[address]; !ok {
		work := &worker{address: address, loggerName: loggerName}
		work.close = make(chan int, 1)
		work.create = make(chan int, 255)
		work.conns = make(chan *conn, 255)
		work.log, _ = logger.Get(loggerName)
		go work.run()
	}
	mutex.Unlock()
	if _, ok := workers[address]; !ok {
		return nil, fmt.Errorf("not find %s", address)
	}
	worker := workers[address]
	if atomic.LoadInt32(&worker.available) == 0 {
		select {
		case worker.create <- 1:
		default:
		}
	}

	timeout := time.NewTicker(time.Second * 3)
	select {
	case c := <-worker.conns:
		atomic.AddInt32(&worker.available, -1)
		return c.client, c.err
	case <-timeout.C:
		return nil, fmt.Errorf("create server %s timeout", address)
	}
}

func removeWorker(address string) {
	mutex.Lock()
	if worker, ok := workers[address]; ok {
		worker.close <- 1
	}
	mutex.Unlock()
}
func removeWorkers() {
	mutex.Lock()
	for _, w := range workers {
		w.close <- 1
	}
	mutex.Unlock()
}

func (w *worker) run() {
START:
	for {
		select {
		case <-w.close:
			break START
		case <-w.create:
			client := NewRPCClientTimeout(w.address, time.Second*5, w.loggerName)
			err := client.Open()
			atomic.AddInt32(&w.available, 1)
			w.conns <- &conn{client: client, err: err}
		}
	}
}
