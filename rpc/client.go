/*
RPC客户端，维护RPC服务器列表，并提供RPC服务调用接口，调用方式分为同步和异步，相同RPC服务有多个服务器时使用轮询机制
选择服务器
该客户端可用于APP->RC, RC->SP ,SP->RC, RC->Job
*/

package rpc

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/colinyl/ars/base"
	"github.com/colinyl/ars/base/rpcservice"
	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/concurrent"
	"github.com/colinyl/lib4go/logger"
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
	queues     concurrent.ConcurrentMap //map[string]chan []interface{}
	pool       *rpcservice.RPCServerPool
	services   concurrent.ConcurrentMap
	client     cluster.IClusterClient
	mutex      sync.RWMutex
	Log        logger.ILogger
	snaps      concurrent.ConcurrentMap
	loggerName string
}

//NewRPCClient 创建RPC Client
func NewRPCClient(cli cluster.IClusterClient, loggerName string) *RPCClient {
	client := &RPCClient{}
	client.client = cli
	client.snaps = concurrent.NewConcurrentMap()
	client.pool = rpcservice.NewRPCServerPool(5, 10, loggerName)
	client.services = concurrent.NewConcurrentMap()
	client.queues = concurrent.NewConcurrentMap()
	client.loggerName = loggerName
	client.Log, _ = logger.Get(loggerName)
	return client
}

//SetPoolSize  设置连接池大小
func (r *RPCClient) SetPoolSize(minSize int, maxSize int) {
	r.pool.MinSize = minSize
	r.pool.MaxSize = maxSize
	r.pool.ResetAllPoolSize(minSize, maxSize)
}

//Close 关闭连接池
func (r *RPCClient) Close() {
	r.pool.Close()
}

func (rc *RPCClient) recover() {
	if r := recover(); r != nil {
		rc.Log.Fatal(r, string(debug.Stack()))
	}
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
				aips = append(aips, ip) //收集IP地址
			}
		}
		if len(v) > 0 {
			r.services.Set(n, &serviceItem{service: v}) //添加新服务
		} else {
			r.services.Delete(n) //移除无可用IP的服务
		}
	}
	for k := range service {
		if _, ok := servers[k]; !ok {
			r.services.Delete(k) //移除已不存在的服务
		}
	}
	r.pool.Register(ips)
	return strings.Join(aips, ",")
}

//GetAsyncResult 获取异步请求结果
func (r *RPCClient) GetAsyncResult(session string) (rt interface{}, err interface{}) {
	queue := r.queues.Get(session)
	if queue != nil {
		ticker := time.NewTicker(time.Second * 4)
		select {
		case <-ticker.C:
			err = fmt.Sprint("request timeout")
			break
		case result := <-queue.(chan []interface{}):
			{
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
			}
		}

	} else {
		err = fmt.Sprint("not find session:", session)
	}
	return
}

//getDomain 获取domain
func (r *RPCClient) getDomain(name string) string {
	if !strings.Contains(name, "@") {
		return ""
	}
	items := strings.Split(name, "@")
	return "@" + items[1]
}

//getGroupName 根据名称获取一个分组
func (r *RPCClient) getGroupName(name string) string {

	group := r.services.Get(name)
	if group == nil {
		group = r.services.Get("*" + r.getDomain(name))
	}
	if group == nil {
		group = r.services.Get("*")
	}

	if group != nil {
		return group.(*serviceItem).getOne()
	}
	return ""
}
func (r *RPCClient) setLifeTime(name string, start time.Time) {
	ss := &base.ProxySnap{}
	ss.ElapsedTime = base.ServerSnap{}
	snap := r.snaps.GetOrAdd(name, ss)
	snap.(*base.ProxySnap).ElapsedTime.Add(start)
}

//Request 发送Request请求
func (r *RPCClient) Request(cmd string, input string, session string) (result string, err error) {
	defer r.recover()
	clogger, _ := logger.NewSession(r.loggerName, session)
	clogger.Info("--> rpc request(send):", cmd, input)
	name := r.client.GetServiceFullPath(cmd)
	group := r.getGroupName(name)
	if strings.EqualFold(group, "") {
		result = base.GetErrorResult("500", "not find rpc server: ", name, " in service list")
		return
	}
	defer r.setLifeTime(group, time.Now())
	result, er := r.pool.Request(group, name, input, session)
	if er != nil {
		result = base.GetErrorResult("500", er.Error())
	} else {
		result = base.GetDataResult(result, false)
	}
	clogger.Info("--> rpc response(send):", cmd, result)
	return
}

//Send 发送Send请求
func (r *RPCClient) Send(cmd string, input string, data string) (result string, err error) {
	name := r.client.GetServiceFullPath(cmd)
	result, _ = r.pool.Send(r.getGroupName(name), name, input, []byte(data))
	return
}

//Get 发送Gety请求
func (r *RPCClient) Get(cmd string, input string) (result string, err error) {
	name := r.client.GetServiceFullPath(cmd)
	data, _ := r.pool.Get(r.getGroupName(name), name, input)
	if err != nil {
		result = string(data)
	}
	return
}

//AsyncRequest 发送异步Request请求
func (r *RPCClient) AsyncRequest(name string, input string, contextSession string) (session string, err error) {
	session = utility.GetGUID()
	queueChan := make(chan []interface{}, 1)
	r.queues.Set(session, queueChan)
	go func(queueChan chan []interface{}, r *RPCClient, name string, input string, csession string) {
		defer r.recover()
		result, err := r.Request(name, input, csession)
		if err != nil {
			queueChan <- []interface{}{result, err.Error()}
		} else {
			queueChan <- []interface{}{result, ""}
		}

	}(queueChan, r, name, input, contextSession)
	return
}

//AsyncSend 发送异步send请求
func (r *RPCClient) AsyncSend(name string, input string, data string) (session string) {
	session = utility.GetGUID()
	queueChan := make(chan []interface{}, 1)
	r.queues.Set(session, queueChan)
	go func(queue chan []interface{}, r *RPCClient, name string, input string, data string) {
		defer r.recover()
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
		defer r.recover()
		result, err := r.Get(name, input)
		if err != nil {
			queue <- []interface{}{result, err.Error()}
		} else {
			queue <- []interface{}{result, ""}
		}

	}(queueChan, r, name, input)
	return
}

//GetSnap 获取RPC客户端的连接池
func (r *RPCClient) GetSnap() []interface{} {
	poolSnaps := r.pool.GetSnap()
	snaps := r.snaps.GetAll()
	return base.GetProxySnap(poolSnaps, snaps)
}
