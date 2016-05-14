/*
RPC客户端，维护RPC服务器列表，并提供RPC服务调用接口，调用方式分为同步和异步，相同RPC服务有多个服务器时使用轮询机制
选择服务器
该客户端可用于APP->RC, RC->SP ,SP->RC, RC->Job
*/

package rpcproxy

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/concurrent"
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
	if len(i.service) == 0 {
		return ""
	}
	index := atomic.AddInt32(&i.index, 1)
	cindex := index % int32(len(i.service))
	return i.service[cindex]
}

//RPCClient RPCClient
type RPCClient struct {
	queues   concurrent.ConcurrentMap //map[string]chan []interface{}
	pool     *rpcservice.RPCServerPool
	services concurrent.ConcurrentMap
	mutex    sync.RWMutex
}

//NewRPCClient 创建RPC Client
func NewRPCClient() *RPCClient {
	client := &RPCClient{}
	client.pool = rpcservice.NewRPCServerPool()
	client.services = concurrent.NewConcurrentMap()
	client.queues = concurrent.NewConcurrentMap()
	return client
}

//SetPoolSize  设置连接池大小
func (r *RPCClient) SetPoolSize(minSize int, maxSize int) {
	r.pool.MinSize = minSize
	r.pool.MaxSize = maxSize
}

//ResetRPCServer 重置所有RPC服务器
func (r *RPCClient) ResetRPCServer(servers map[string][]string) string {
	ips := make(map[string]string) //构建IP列表，用于注册服务
	aips := []string{}
	service := r.services.GetAll()
	for n, v := range servers {
		for _, ip := range v {
			if _, ok := ips[ip]; !ok {
				ips[ip] = ip
				aips = append(aips, ip)
			}
		}
		if _, ok := service[n]; !ok {
			r.services.Set(n, &serviceItem{service: v})
		}
	}
	r.pool.Register(ips)
	return strings.Join(aips, ",")
}

//GetAsyncResult 获取异步请求结果
func (r *RPCClient) GetAsyncResult(session string) (rt interface{}, err interface{}) {
	queue := r.queues.Get(session)
	if queue != nil {
		result := <-queue.(chan []interface{})
		r.queues.Delete(session)
		if len(result) != 2 {
			return "", "rpc method result value len is error"
		}
		rt = result[0]
		if result[1] != nil {
			er := result[1].(string)
			if strings.EqualFold(er, "") {
				err = nil
			} else {
				err = er
			}
		}
	} else {
		err = fmt.Sprint("not find session:", session)
	}
	return
}

//getGroupName 根据名称获取一个分组
func (r *RPCClient) getGroupName(name string) string {
	group := r.services.Get(name)
	if group == nil {
		group = r.services.Get("*")
	}
	if group != nil {
		return group.(*serviceItem).getOne()
	}
	return ""
}

//Request 发送Request请求
func (r *RPCClient) Request(name string, input string) (result string, err error) {
	group := r.getGroupName(name)
	if strings.EqualFold(group, "") {
		return "", fmt.Errorf("not find rpc server:%s", group)
	}
	result, er := r.pool.Request(group, name, input)
	if er != nil {
		result = GetErrorResult("500", er.Error())
	} else {
		result = GetDataResult(result)
	}
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
	r.queues.Set(session, queueChan)
	go func(queueChan chan []interface{}, r *RPCClient, name string, input string) {
		result, err := r.Request(name, input)
		if err != nil {
			queueChan <- []interface{}{result, err.Error()}
		} else {
			queueChan <- []interface{}{result, ""}
		}

	}(queueChan, r, name, input)
	return
}

//AsyncSend 发送异步send请求
func (r *RPCClient) AsyncSend(name string, input string, data string) (session string) {
	session = utility.GetGUID()
	queueChan := make(chan []interface{}, 1)
	r.queues.Set(session, queueChan)
	go func(queue chan []interface{}, r *RPCClient, name string, input string, data string) {
		result, err := r.Send(name, input, data)
		if err != nil {
			queue <- []interface{}{result, err.Error()}
		} else {
			queue <- []interface{}{result, ""}
		}

	}(queueChan, r, name, input, data)
	return
}

//AsyncGet 发送异步GET请求
func (r *RPCClient) AsyncGet(name string, input string) (session string) {
	session = utility.GetGUID()
	queueChan := make(chan []interface{}, 1)
	r.queues.Set(session, queueChan)
	go func(queue chan []interface{}, r *RPCClient, name string, input string) {
		result, err := r.Get(name, input)
		if err != nil {
			queue <- []interface{}{result, err.Error()}
		} else {
			queue <- []interface{}{result, ""}
		}

	}(queueChan, r, name, input)
	return
}
