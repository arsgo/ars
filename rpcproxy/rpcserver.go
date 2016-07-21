/*RPC服务器，动态监听本机端口，需外部指定RPC服务器处理程序(RPCHandler),处理程序目前已实现的包括:
RPC转发处理程序(RPCProxyHandler),用于RC服务器转发APP请求
Script执行程序(RPCScriptHandler),用于执行本地LUA脚本
*/
package rpcproxy

import (
	"fmt"
	"strings"
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
	Log           logger.ILogger
	loggerName    string
	snap          *ServerSnap
}

//Tasks 任务列表
type Tasks struct {
	Services concurrent.ConcurrentMap
}

//RPCHandler RPC处理函数
type RPCHandler interface {
	OpenTask(cluster.TaskItem)
	CloseTask(cluster.TaskItem)
	Request(cluster.TaskItem, string, string) (string, error)
	Send(cluster.TaskItem, string, []byte) (string, error)
	Get(cluster.TaskItem, string) ([]byte, error)
}

//NewRPCServer 创建RPC服务器
func NewRPCServer(handler RPCHandler, loggerName string) (server *RPCServer) {
	server = &RPCServer{loggerName: loggerName, snap: &ServerSnap{}}
	server.serverHandler = NewRPCHandlerProxy(handler, loggerName, server.snap)
	server.Log, _ = logger.Get(loggerName, true)
	return server
}

//UpdateTasks 更新任务列表
func (r *RPCServer) UpdateTasks(tasks []cluster.TaskItem) {
	r.serverHandler.UpdateTasks(tasks)
}

//Start 启动RPC服务器
func (r *RPCServer) Start() {
	r.Address = rpcservice.GetLocalRandomAddress()
	r.server = rpcservice.NewRPCServer(r.Address, r.serverHandler, r.loggerName)
	r.server.Serve()
	time.Sleep(time.Second)
}

//Stop 停止RPC服务
func (r *RPCServer) Stop() {
	if r.server != nil {
		r.server.Stop()
	}
}

//GetSnap 获取当前服务器快照信息
func (r *RPCServer) GetSnap() ServerSnap {
	return *r.snap
}

//RPCHandlerProxy RPCHandler代理程序
type RPCHandlerProxy struct {
	tasks      Tasks
	handler    RPCHandler
	Log        logger.ILogger
	snap       *ServerSnap
	loggerName string
}

//NewRPCHandlerProxy 创建RPC默认处理程序
func NewRPCHandlerProxy(h RPCHandler, loggerName string, snap *ServerSnap) *RPCHandlerProxy {
	handler := &RPCHandlerProxy{snap: snap, loggerName: loggerName}
	handler.tasks = Tasks{}
	handler.handler = h
	handler.tasks.Services = concurrent.NewConcurrentMap() //make(map[string]cluster.TaskItem)
	handler.Log, _ = logger.Get(loggerName, true)
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
			r.handler.CloseTask(v.(cluster.TaskItem))
			r.tasks.Services.Delete(i)
		} else {
			r.tasks.Services.Set(i, tks[i]) //更新可能已经变化的服务
		}
	}
	for i, v := range tks {
		if _, ok := services[i]; !ok {
			r.tasks.Services.Set(i, v) //添加新任务
			r.handler.OpenTask(v)
		}
	}
}

//getDomain 获取domain
func (r *RPCHandlerProxy) getDomain(name string) string {
	if !strings.Contains(name, "@") {
		return ""
	}
	items := strings.Split(name, "@")
	return "@" + items[1]
}

//getTaskItem 根据名称获取一个分组
func (r *RPCHandlerProxy) getTaskItem(name string) (item cluster.TaskItem, err error) {

	//	all := r.tasks.Services.GetAll()
	//	for i, v := range all {
	//fmt.Printf("getTaskItem:%s,%v\n", i, v.(cluster.TaskItem).IP)
	//	}

	//r.Log.Info("get1:", name)
	group := r.tasks.Services.Get(name)
	if group == nil {
		//r.Log.Info("get3:", "*"+r.getDomain(name))
		group = r.tasks.Services.Get("*" + r.getDomain(name))
	}
	if group == nil {
		//r.Log.Info("get3:", "*")
		group = r.tasks.Services.Get("*")
	}

	if group != nil {
		item = group.(cluster.TaskItem)
		item.Name = name
		return
	}
	err = fmt.Errorf("not find service:%s", name)
	return
}

//Request 执行RPC Request服务
func (r *RPCHandlerProxy) Request(name string, input string, session string) (result string, err error) {
	defer r.snap.Add(time.Now())
	log, _ := logger.NewSession(r.loggerName, session, true)
	log.Info("--> rpc request:", name, input)
	task, currentErr := r.getTaskItem(name)
	if currentErr != nil {
		result = GetErrorResult("500", currentErr.Error())
	} else {
		result, currentErr = r.handler.Request(task, input, session)
	}
	if currentErr != nil {
		r.Log.Error(currentErr)
	}
	log.Infof("--> rpc response:", name, result)
	return
}

//Send 执行RPC Send服务
func (r *RPCHandlerProxy) Send(name string, input string, data []byte) (result string, err error) {
	r.Log.Info("-> recv send:", name)
	task, er := r.getTaskItem(name)
	if er != nil {
		return GetErrorResult("500", er.Error()), nil
	}
	result, er = r.handler.Send(task, input, data)
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
	task, er := r.getTaskItem(name)
	if er != nil {
		return []byte(GetErrorResult("500", er.Error())), nil
	}
	buffer, er = r.handler.Get(task, input)
	if er != nil {
		r.Log.Error(er)
		buffer = []byte(GetErrorResult("500", er.Error()))
	}
	return
}
