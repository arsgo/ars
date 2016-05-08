package rpcserver

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/rpcservice"
)

//RPCServer RPC服务器
type RPCServer struct {
	Address       string
	serverHandler *rpcHandler
	server        *rpcservice.RPCServer
}

//Tasks 任务列表
type Tasks struct {
	Services map[string]cluster.TaskItem
	mutex    sync.RWMutex
}

//RPCScriptHandler RPC处理函数
type RPCScriptHandler interface {
	OpenTask(cluster.TaskItem)
	CloseTask(cluster.TaskItem)
	Request(cluster.TaskItem, string) (string, error)
	Send(cluster.TaskItem, string, []byte) (string, error)
	Get(cluster.TaskItem, string) ([]byte, error)
}

//NewRPCServer 创建RPC服务器
func NewRPCServer(handler RPCScriptHandler) (server *RPCServer) {
	server = &RPCServer{}
	server.serverHandler = NewRPCHandler(handler)
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

type rpcHandler struct {
	tasks   Tasks
	handler RPCScriptHandler
}

//NewRPCHandler 创建RPC默认处理程序
func NewRPCHandler(h RPCScriptHandler) *rpcHandler {
	handler := &rpcHandler{}
	handler.tasks = Tasks{}
	handler.handler = h
	handler.tasks.Services = make(map[string]cluster.TaskItem)
	return handler
}

//UpdateTasks 更新服务列表
func (r *rpcHandler) UpdateTasks(tasks []cluster.TaskItem) {
	tks := make(map[string]cluster.TaskItem)
	for _, v := range tasks {
		tks[v.Name] = v
	}
	r.tasks.mutex.Lock()
	for i, v := range r.tasks.Services {
		if _, ok := tks[i]; !ok {
			r.handler.CloseTask(v)
			delete(r.tasks.Services, i)
		}
	}
	for i, v := range tks {
		if _, ok := r.tasks.Services[i]; !ok {
			r.tasks.Services[i] = v
			r.handler.OpenTask(v)
		}
	}
	r.tasks.mutex.Unlock()
}

//Request 执行RPC Request服务
func (r *rpcHandler) Request(name string, input string) (result string, err error) {
	r.tasks.mutex.RLock()
	task, ok := r.tasks.Services[name]
	r.tasks.mutex.Unlock()
	if !ok || strings.EqualFold(strings.ToLower(task.Method), "request") {
		err = errors.New(fmt.Sprint("not find rpc service:", name))
		return
	}
	result, err = r.handler.Request(task, input)
	return
}

//Send 执行RPC Send服务
func (r *rpcHandler) Send(name string, input string, data []byte) (result string, err error) {
	r.tasks.mutex.RLock()
	task, ok := r.tasks.Services[name]
	r.tasks.mutex.Unlock()
	if !ok || strings.EqualFold(strings.ToLower(task.Method), "send") {
		err = errors.New(fmt.Sprint("not find rpc service:", name))
		return
	}
	return r.handler.Send(task, input, data)
}

//Get 执行RPC Get服务
func (r *rpcHandler) Get(name string, input string) (buffer []byte, err error) {
	r.tasks.mutex.RLock()
	task, ok := r.tasks.Services[name]
	r.tasks.mutex.Unlock()
	if !ok || strings.EqualFold(strings.ToLower(task.Method), "get") {
		err = errors.New(fmt.Sprint("not find rpc service:", name))
		return
	}
	return r.handler.Get(task, input)
}
