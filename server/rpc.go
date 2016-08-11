/*RPC服务器，动态监听本机端口，需外部指定RPC服务器处理程序(RPCHandler),处理程序目前已实现的包括:
RPC转发处理程序(RPCProxyHandler),用于RC服务器转发APP请求
Script执行程序(RPCScriptHandler),用于执行本地LUA脚本
*/
package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/base/rpcservice"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/ars/servers/config"
	"github.com/arsgo/lib4go/concurrent"
	"github.com/arsgo/lib4go/logger"
)

//RPCServer RPC服务器
type RPCServer struct {
	Address       string
	serverHandler *RPCHandlerProxy
	server        *rpcservice.RPCServer
	Log           logger.ILogger
	loggerName    string
	snap          *base.ServerSnap
	collector     *base.Collector
}

//Tasks 任务列表
type Tasks struct {
	Services *concurrent.ConcurrentMap
}

//IRPCHandler RPC处理函数
type IRPCHandler interface {
	OpenTask(cluster.TaskItem)
	CloseTask(cluster.TaskItem)
	Request(cluster.TaskItem, string, string) (string, error)
	Send(cluster.TaskItem, string, []byte) (string, error)
	Get(cluster.TaskItem, string) ([]byte, error)
}

//NewRPCServer 创建RPC服务器
func NewRPCServer(handler IRPCHandler, loggerName string, callback base.CollectorCallBack) (server *RPCServer) {
	server = &RPCServer{loggerName: loggerName, snap: &base.ServerSnap{}}
	server.collector = base.NewCollector(callback, time.Second)
	server.serverHandler = NewRPCHandlerProxy(handler, loggerName, server.snap, server.collector)
	server.Log, _ = logger.Get(loggerName)
	return server
}

//UpdateTasks 更新任务列表
func (r *RPCServer) UpdateTasks(tasks []cluster.TaskItem) int {
	return r.serverHandler.UpdateTasks(tasks)
}

//GetServices 获取所有服务信息
func (r *RPCServer) GetServices() []string {
	return r.serverHandler.GetServices()
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
func (r *RPCServer) GetSnap() base.ServerSnap {
	return *r.snap
}

//RPCHandlerProxy RPCHandler代理程序
type RPCHandlerProxy struct {
	tasks      Tasks
	handler    IRPCHandler
	Log        logger.ILogger
	snap       *base.ServerSnap
	collector  *base.Collector
	domain     string
	loggerName string
}

//NewRPCHandlerProxy 创建RPC默认处理程序
func NewRPCHandlerProxy(h IRPCHandler, loggerName string, snap *base.ServerSnap, collector *base.Collector) *RPCHandlerProxy {
	conf, _ := config.Get()
	handler := &RPCHandlerProxy{snap: snap, loggerName: loggerName, domain: conf.Domain, collector: collector}
	handler.tasks = Tasks{}
	handler.handler = h
	handler.tasks.Services = concurrent.NewConcurrentMap() //make(map[string]cluster.TaskItem)
	handler.Log, _ = logger.Get(loggerName)
	return handler
}

//GetServices 获取所有服务信息
func (r *RPCHandlerProxy) GetServices() (v []string) {
	svs := r.tasks.Services.GetAll()
	v = make([]string, 0, len(svs))
	for i := range svs {
		v = append(v, i)
	}
	return v
}

//UpdateTasks 更新服务列表
func (r *RPCHandlerProxy) UpdateTasks(tasks []cluster.TaskItem) int {
	count := 0
	tks := make(map[string]cluster.TaskItem)
	for _, v := range tasks {
		tks[v.Name] = v
	}
	services := r.tasks.Services.GetAll()
	for i, v := range tks {
		if _, ok := services[i]; !ok {
			if r.tasks.Services.Set(i, v) { //添加新任务
				count++
			}
			r.handler.OpenTask(v)
		}
	}
	for i, v := range services {
		if _, ok := tks[i]; !ok {
			r.handler.CloseTask(v.(cluster.TaskItem))
			r.tasks.Services.Delete(i)
		} else {
			if r.tasks.Services.Set(i, tks[i]) { //更新可能已经变化的服务
				count++
			}
		}
	}
	return count
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
	group := r.tasks.Services.Get(name)
	if group == nil {
		group = r.tasks.Services.Get("*" + r.getDomain(name))
	}
	if group == nil {
		group = r.tasks.Services.Get("*")
	}
	if group != nil {
		item = group.(cluster.TaskItem)
		item.Name = name
		return
	}
	r.collector.Error()
	err = fmt.Errorf("not find service(%s@%s.rpc.server):%s,%d", r.loggerName, r.domain, name, r.tasks.Services.GetLength())
	return
}

//Request 执行RPC Request服务
func (r *RPCHandlerProxy) Request(name string, input string, session string) (result string, err error) {
	defer r.snap.Add(time.Now())
	log, _ := logger.NewSession(r.loggerName, session)
	log.Info("--> rpc request(recv):", name, input)
	task, currentErr := r.getTaskItem(name)
	if currentErr != nil {
		result = base.GetErrorResult("500", currentErr.Error())
	} else {
		result, currentErr = r.handler.Request(task, input, session)
	}
	if currentErr != nil {
		r.Log.Error(currentErr)
		r.collector.Failed()
	} else {
		r.collector.Success()
	}
	log.Info("--> rpc response(recv):", name, result)
	return
}

//Send 执行RPC Send服务
func (r *RPCHandlerProxy) Send(name string, input string, data []byte) (result string, err error) {
	r.Log.Info("-> recv send:", name)
	task, er := r.getTaskItem(name)
	if er != nil {
		return base.GetErrorResult("500", er.Error()), nil
	}
	result, er = r.handler.Send(task, input, data)
	if er != nil {
		r.Log.Error(er)
		result = base.GetErrorResult("500", er.Error())
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
		return []byte(base.GetErrorResult("500", er.Error())), nil
	}
	buffer, er = r.handler.Get(task, input)
	if er != nil {
		r.Log.Error(er)
		buffer = []byte(base.GetErrorResult("500", er.Error()))
	}
	return
}
