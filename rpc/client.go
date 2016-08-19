/*
RPC客户端，维护RPC服务器列表，并提供RPC服务调用接口，调用方式分为同步和异步，相同RPC服务有多个服务器时使用轮询机制
选择服务器
该客户端可用于APP->RC, RC->SP ,SP->RC, RC->Job
*/

package rpc

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/base/rpcservice"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/ars/servers/config"
	"github.com/arsgo/lib4go/concurrent"
	"github.com/arsgo/lib4go/logger"
	"github.com/arsgo/lib4go/utility"
)

//RPCClient RPCClient
type RPCClient struct {
	queues     *concurrent.ConcurrentMap //map[string]chan []interface{}
	pool       *rpcservice.RPCServerPool
	services   *concurrent.ConcurrentMap
	client     cluster.IClusterClient
	mutex      sync.RWMutex
	Log        logger.ILogger
	snaps      *concurrent.ConcurrentMap
	loggerName string
	domain     string
}

//NewRPCClient 创建RPC Client
func NewRPCClient(cli cluster.IClusterClient, loggerName string) *RPCClient {
	conf, _ := config.Get()
	client := &RPCClient{domain: conf.Domain}
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

func (r *RPCClient) GetServices() map[string]interface{} {
	return r.services.GetAll()
}

func (r *RPCClient) GetServiceCount() int {
	return r.services.GetLength()
}

//ResetRPCServer 重置所有RPC服务器
func (r *RPCClient) ResetRPCServer(servers map[string][]string) string {
	ips := make(map[string]string) //构建IP列表，用于注册服务
	aips := []string{}
	service := r.services.GetAll()
	var setServer, delServer []string
	for n, v := range servers {
		for _, ip := range v {
			if _, ok := ips[ip]; !ok {
				ips[ip] = ip
				aips = append(aips, ip) //收集IP地址
			}
		}
		if len(v) > 0 {
			setServer = append(setServer, n)
			r.services.Set(n, base.NewServiceItem(v)) //添加新服务
		} else {
			delServer = append(delServer, n)
			r.services.Delete(n) //移除无可用IP的服务
		}
	}
	for k := range service {
		if _, ok := servers[k]; !ok {
			delServer = append(delServer, k)
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

//getGroup 根据名称获取一个分组
func (r *RPCClient) getGroup(cmd string) (g *base.ServiceGroup, name string, err error) {
	name = r.client.GetServiceFullPath(cmd)
	group := r.services.Get(name)
	if group == nil {
		group = r.services.Get("*" + r.getDomain(name))
	}
	if group == nil {
		group = r.services.Get("*")
	}
	if group != nil {
		g = group.(*base.ServiceItem).GetGroup()
	} else {
		err = errors.New(fmt.Sprint("not find rpc server(", r.loggerName, "@", r.domain, ".rpc.client):", name, " in service list",
			r.services.GetLength()))
	}
	return
}

//Request 发送Request请求
func (r *RPCClient) Request(cmd string, input string, session string) (result string, err error) {
	defer r.recover()
	start := time.Now()
	clogger, _ := logger.NewSession(r.loggerName, session)
	clogger.Info("--> rpc request(send):", cmd, input)
	group, name, err := r.getGroup(cmd)
	if err != nil {
		result = base.GetErrorResult(base.ERR_NOT_FIND_SRVS, err.Error())
		return
	}
	groupName := cmd
	defer func() {
		clogger.Infof("--> rpc response(send,%v):%s,%s", time.Now().Sub(start), cmd, result)
	}()
	defer r.setLifeTime(groupName, time.Now())
START:
	groupName, err = group.GetNext()
	if err != nil {
		if strings.EqualFold(result, "") {
			msg := fmt.Sprint("not find rpc server(", r.loggerName, "@", r.domain, ".rpc.client):", name, " in service list",
				r.services.GetLength())
			result = base.GetErrorResult(base.ERR_NOT_FIND_SRVS, msg)
		}
		err = errors.New(result)
		return
	}
	result, err = r.pool.Request(groupName, name, input, session)
	if err != nil {
		result = base.GetErrorResult(base.ERR_NOT_FIND_SRVS, err.Error())
	} else {
		result = base.GetDataResult(result, false)
	}
	if strings.EqualFold(base.GetResult(result).Code, base.ERR_NOT_FIND_SRVS) {
		goto START
	}

	return
}

//Send 发送Send请求
func (r *RPCClient) Send(cmd string, input string, data string) (result string, err error) {
	return
}

//Get 发送Gety请求
func (r *RPCClient) Get(cmd string, input string) (result string, err error) {
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
func (r *RPCClient) createSnap(p ...interface{}) (interface{}, error) {
	ss := &base.ProxySnap{}
	ss.ElapsedTime = base.ServerSnap{}
	return ss, nil
}
func (r *RPCClient) setLifeTime(name string, start time.Time) {
	_, snap, _ := r.snaps.Add(name, r.createSnap)
	if snap == nil {
		return
	}
	snap.(*base.ProxySnap).ElapsedTime.Add(start)
}
