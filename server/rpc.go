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
	servicesPath  *concurrent.ConcurrentMap
	Address       string
	Available     bool
	serverHandler *RPCHandlerProxy
	server        *rpcservice.RPCServer
	Log           logger.ILogger
	loggerName    string
	snap          *base.ServerSnap
	collector     base.ICollector
}

//Tasks 任务列表
type Tasks struct {
	Services *concurrent.ConcurrentMap
}

//IRPCHandler RPC处理函数
type IRPCHandler interface {
	OpenTask(cluster.TaskItem) string
	CloseTask(cluster.TaskItem, string)
	Request(cluster.TaskItem, string, string) (string, error)
	Send(cluster.TaskItem, string, []byte) (string, error)
	Get(cluster.TaskItem, string) ([]byte, error)
}

//NewRPCServer 创建RPC服务器
func NewRPCServer(handler IRPCHandler, loggerName string, collector base.ICollector) (server *RPCServer) {
	server = &RPCServer{loggerName: loggerName, snap: &base.ServerSnap{}, collector: collector}
	server.serverHandler = NewRPCHandlerProxy(server, handler, loggerName, server.snap, server.collector)
	server.servicesPath = concurrent.NewConcurrentMap()
	server.Log, _ = logger.Get(loggerName)
	return server
}

//UpdateTasks 更新任务列表
func (r *RPCServer) UpdateTasks(tasks []cluster.TaskItem) (c int) {
	c, r.Available = r.serverHandler.UpdateTasks(tasks)
	return
}

//GetServicePath  获取当前服务路径
func (r *RPCServer) GetServicePath() (paths map[string]string) {
	paths = make(map[string]string)
	services := r.servicesPath.GetAll()
	for name, path := range services {
		paths[name] = path.(string)
	}
	return
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
	server     *RPCServer
	tasks      Tasks
	handler    IRPCHandler
	Available  bool
	Log        logger.ILogger
	snap       *base.ServerSnap
	collector  base.ICollector
	domain     string
	loggerName string
}

//NewRPCHandlerProxy 创建RPC默认处理程序
func NewRPCHandlerProxy(server *RPCServer, h IRPCHandler, loggerName string, snap *base.ServerSnap, collector base.ICollector) *RPCHandlerProxy {
	conf, _ := config.Get()
	handler := &RPCHandlerProxy{server: server, snap: snap, loggerName: loggerName, domain: conf.Domain, collector: collector, Available: false}
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
func (r *RPCHandlerProxy) UpdateTasks(tasks []cluster.TaskItem) (int, bool) {
	count := 0
	tks := make(map[string]cluster.TaskItem)
	for _, v := range tasks {
		tks[v.Name] = v
	}
	services := r.tasks.Services.GetAll()
	for i, v := range tks {
		if _, ok := services[i]; !ok {
			if r.tasks.Services.Set(i, v) { //添加新任务
				path := r.handler.OpenTask(v)
				r.server.servicesPath.Set(i, path)
				count++
			}
		}
	}
	for i, v := range services {
		if _, ok := tks[i]; !ok {
			tk := v.(cluster.TaskItem)
			value, ok := r.server.servicesPath.Get(tk.Name)
			if ok {
				r.handler.CloseTask(tk, value.(string))
			}
			r.tasks.Services.Delete(i)
			count++
		} else {
			if r.tasks.Services.Set(i, tks[i]) { //更新可能已经变化的服务
				count++
			}
		}
	}
	r.Available = r.tasks.Services.GetLength() > 0
	return count, r.Available
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
	group, ok := r.tasks.Services.Get(name)
	if !ok {
		group, ok = r.tasks.Services.Get("*" + r.getDomain(name))
	}
	if !ok {
		group, ok = r.tasks.Services.Get("*")
	}
	if ok {
		item = group.(cluster.TaskItem)
		item.Name = name
		return
	}
	err = fmt.Errorf("not find service(%s@%s.rpc.server):%s,%d", r.loggerName, r.domain, name, r.tasks.Services.GetLength())
	return
}

//Request 执行RPC Request服务
func (r *RPCHandlerProxy) Request(name string, input string, session string) (result string, err error) {
	defer r.snap.Add(time.Now())
	start := time.Now()
	//	defer base.RunTime("rpc request once", time.Now())
	log, _ := logger.NewSession(r.loggerName, session)
	log.Info("--> rpc request(recv):", name, input)

	task, currentErr := r.getTaskItem(name)
	if currentErr != nil {
		result = base.GetErrorResult(base.ERR_NOT_FIND_SRVS, currentErr.Error())
		r.Log.Error(currentErr)
		r.collector.Error(name)
		log.Infof("--> rpc response(recv,%v):%s,%s", time.Now().Sub(start), name, result)
		return
	}
	result, currentErr = r.handler.Request(task, input, session)
	if currentErr != nil {
		r.Log.Error(currentErr)
		r.collector.Failed(name)
		log.Infof("--> rpc response(recv,%v):%s,%s", time.Now().Sub(start), name, result)
		return
	}
	if base.GetResult(result).Code == base.ERR_NOT_FIND_SRVS {
		r.collector.Error(name)
		log.Infof("--> rpc response(recv,%v):%s,%s", time.Now().Sub(start), name, result)
		return
	}
	r.collector.Success(name)
	log.Infof("--> rpc response(recv,%v):%s,%s", time.Now().Sub(start), name, result)
	return
}

//Send 执行RPC Send服务
func (r *RPCHandlerProxy) Send(name string, input string, data []byte) (result string, err error) {
	return
}

//Heartbeat 返回心跳数据
func (r *RPCHandlerProxy) Heartbeat(input string) (rs string, err error) {
	return "success", nil
}

//Get 执行RPC Get服务
func (r *RPCHandlerProxy) Get(name string, input string) (buffer []byte, err error) {

	return
}
