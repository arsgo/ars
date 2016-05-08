package rpcclient

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/utility"
)

//serviceItem 服务信息
type serviceItem struct {
	service []string
	index   int32
	mutex   sync.Mutex
}

//getOne 获取一个可用的服务
func (i *serviceItem) getOne() string {
	index := atomic.AddInt32(&i.index, 1)
	cindex := index % int32(len(i.service))
	return i.service[cindex]
}

//RPCClient RPCClient
type RPCClient struct {
	queues   map[string]chan []interface{}
	pool     *rpcservice.RPCServerPool
	services map[string]serviceItem
	mutex    sync.RWMutex
}

//NewRPCClient 创建RPC Client
func NewRPCClient() *RPCClient {
	client := &RPCClient{}
	client.pool = rpcservice.NewRPCServerPool()
	client.services = make(map[string]serviceItem)
	client.queues = make(map[string]chan []interface{})
	return client
}

//ResetRPCServer 重置所有RPC服务器
func (r *RPCClient) ResetRPCServer(servers map[string][]string) {
	ips := make(map[string]string)
	svs := make(map[string]serviceItem)
	for n, v := range servers {
		for _, ip := range v {
			if _, ok := ips[ip]; !ok {
				ips[ip] = ip
			}
		}
		if _, ok := r.services[n]; !ok {
			svs[n] = serviceItem{service: v}
		}
	}
	r.pool.Register(ips)
	r.mutex.Lock()
	r.services = svs
	r.mutex.Unlock()
}

//GetAsyncResult 获取异步请求结果
func (r *RPCClient) GetAsyncResult(session string) (rt interface{}, err interface{}) {
	r.mutex.RLock()
	if queue, ok := r.queues[session]; ok {
		result := <-queue
		defer delete(r.queues, session)
		if len(result) != 2 {
			return "", errors.New("rpc method result value len is error")
		}
		rt = result[0]
		err = result[1]
	} else {
		err = errors.New(fmt.Sprint("not find session:", session))
	}
	r.mutex.RUnlock()
	return
}

//getGroupName 根据名称获取一个分组
func (r *RPCClient) getGroupName(name string) string {
	r.mutex.RLock()
	group, ok := r.services[name]
	r.mutex.RUnlock()
	if ok {
		return group.getOne()
	}
	group, ok = r.services["*"]
	if ok {
		return group.getOne()
	}
	return ""
}

//Request 发送Request请求
func (r *RPCClient) Request(name string, input string) (result string, err error) {
	group := r.getGroupName(name)
	if strings.EqualFold(group, "") {
		return "", errors.New("not find rpc server")
	}
	result, err = r.pool.Request(group, name, input)
	return
}

//Send 发送Send请求
func (r *RPCClient) Send(name string, input string, data string) (result string, err error) {
	result, err = r.pool.Send(r.getGroupName(name), name, input, []byte(data))
	return
}

//Get 发送Gety请求
func (r *RPCClient) Get(name string, input string) (result string, err error) {
	data, err := r.pool.Get(r.getGroupName(name), name, input)
	if err != nil {
		result = string(data)
	}
	return
}

//AsyncRequest 发送异步Request请求
func (r *RPCClient) AsyncRequest(name string, input string) (session string) {
	session = utility.GetGUID()
	queueChan := make(chan []interface{}, 1)
	r.mutex.Lock()
	r.queues[session] = queueChan
	r.mutex.Unlock()
	go func(queueChan chan []interface{}, r *RPCClient, name string, input string) {
		result, err := r.Request(name, input)
		if err != nil {
			queueChan <- []interface{}{result, err.Error()}
		} else {
			queueChan <- []interface{}{result, nil}
		}

	}(queueChan, r, name, input)
	return
}

//AsyncSend 发送异步send请求
func (r *RPCClient) AsyncSend(name string, input string, data string) (session string) {
	session = utility.GetGUID()
	queueChan := make(chan []interface{}, 1)
	r.mutex.Lock()
	r.queues[session] = queueChan
	r.mutex.Unlock()
	go func(queue chan []interface{}, r *RPCClient, name string, input string, data string) {
		result, err := r.Send(name, input, data)
		if err != nil {
			queue <- []interface{}{result, err.Error()}
		} else {
			queue <- []interface{}{result, nil}
		}

	}(queueChan, r, name, input, data)
	return
}

//AsyncGet 发送异步GET请求
func (r *RPCClient) AsyncGet(name string, input string) (session string) {
	session = utility.GetGUID()
	queueChan := make(chan []interface{}, 1)
	r.mutex.Lock()
	r.queues[session] = queueChan
	r.mutex.Unlock()
	go func(queue chan []interface{}, r *RPCClient, name string, input string) {
		result, err := r.Get(name, input)
		if err != nil {
			queue <- []interface{}{result, err.Error()}
		} else {
			queue <- []interface{}{result, nil}
		}

	}(queueChan, r, name, input)
	return
}
