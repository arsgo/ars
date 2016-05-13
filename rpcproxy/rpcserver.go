/*RPC服务器，动态监听本机端口，需外部指定RPC服务器处理程序(RPCHandler),处理程序目前已实现的包括:
RPC转发处理程序(RPCProxyHandler),用于RC服务器转发APP请求
Script执行程序(RPCScriptHandler),用于执行本地LUA脚本
*/
package rpcproxy

import (
	"time"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/concurrent"
	"github.com/colinyl/lib4go/logger"
)

//RPCServer RPC服务器
type RPCServer struct {
	Address       string
	serverHandler *RPCHandlerProxy
	server        *rpcservice.RPCServer
	Log           *logger.Logger
}

//Tasks 任务列表
type Tasks struct {
	Services concurrent.ConcurrentMap
}

//RPCHandler RPC处理函数
type RPCHandler interface {
	OpenTask(cluster.TaskItem)
	CloseTask(cluster.TaskItem)
	Request(cluster.TaskItem, string) (string, error)
	Send(cluster.TaskItem, string, []byte) (string, error)
	Get(cluster.TaskItem, string) ([]byte, error)
}

//NewRPCServer 创建RPC服务器
func NewRPCServer(handler RPCHandler) (server *RPCServer) {
	server = &RPCServer{}
	server.serverHandler = NewRPCHandlerProxy(handler)
	server.Log, _ = logger.New("rpc server", true)
	return server
}

//UpdateTasks 更新任务列表
func (r *RPCServer) UpdateTasks(tasks []cluster.TaskItem) {
	r.serverHandler.UpdateTasks(tasks)
}

//Start 启动RPC服务器
func (r *RPCServer) Start() {
	r.Address = rpcservice.GetLocalRandomAddress()
	r.server = rpcservice.NewRPCServer(r.Address, r.serverHandler)
	r.server.Serve()
	time.Sleep(time.Second * 2)
}

//Stop 停止RPC服务
func (r *RPCServer) Stop() {
	if r.server != nil {
		r.server.Stop()
	}
}

//RPCHandlerProxy RPCHandler代理程序
type RPCHandlerProxy struct {
	tasks   Tasks
	handler RPCHandler
	Log     *logger.Logger
}

//NewRPCHandlerProxy 创建RPC默认处理程序
func NewRPCHandlerProxy(h RPCHandler) *RPCHandlerProxy {
	handler := &RPCHandlerProxy{}
	handler.tasks = Tasks{}
	handler.handler = h
	handler.tasks.Services = concurrent.NewConcurrentMap() //make(map[string]cluster.TaskItem)
	handler.Log, _ = logger.New("rpc server", true)
	return handler
}

//UpdateTasks 更新服务列表
func (r *RPCHandlerProxy) UpdateTasks(tasks []cluster.TaskItem) {
	tks := make(map[string]cluster.TaskItem)
	for _, v := range tasks {
		tks[v.Name] = v
	}
	services := r.tasks.Services.GetAll()

	for i, v := range services {
		if _, ok := tks[i]; !ok {
			r.tasks.Services.Delete(i)
			go r.handler.CloseTask(v.(cluster.TaskItem))
		} else {
			r.tasks.Services.Set(i, tks[i]) //更新可能已经变化的服务
		}
	}
	for i, v := range tks {
		if _, ok := services[i]; !ok {
			r.tasks.Services.Set(i, v) //添加新任务
			go r.handler.OpenTask(v)
		}
	}

}

//Request 执行RPC Request服务
func (r *RPCHandlerProxy) Request(name string, input string) (result string, err error) {
	r.Log.Info("-> recv request:", name)
	task := r.tasks.Services.Get(name)
	if task != nil {
		return GetErrorResult("500", "not find service:", name), nil
	}
	result, er := r.handler.Request(task.(cluster.TaskItem), input)
	if er != nil {
		r.Log.Error(er)
		result = GetErrorResult("500", er.Error())
	} else {
		r.Log.Info(result)
	}
	return
}

//Send 执行RPC Send服务
func (r *RPCHandlerProxy) Send(name string, input string, data []byte) (result string, err error) {
	r.Log.Info("-> recv send:", name)
	task := r.tasks.Services.Get(name)
	if task == nil {
		return GetErrorResult("500", "not find service:", name), nil
	}
	result, er := r.handler.Send(task.(cluster.TaskItem), input, data)
	if er != nil {
		r.Log.Error(er)
		result = GetErrorResult("500", er.Error())
	} else {
		r.Log.Info(result)
	}
	return
}

//Get 执行RPC Get服务
func (r *RPCHandlerProxy) Get(name string, input string) (buffer []byte, err error) {
	r.Log.Info("-> recv get:", name)
	task := r.tasks.Services.Get(name)
	if task != nil {
		return []byte(GetErrorResult("500", "not find service:", name)), nil
	}
	buffer, er := r.handler.Get(task.(cluster.TaskItem), input)
	if er != nil {
		r.Log.Error(er)
		buffer = []byte(GetErrorResult("500", er.Error()))
	}
	return
}
