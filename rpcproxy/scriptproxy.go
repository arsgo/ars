package rpcproxy

import (
	"errors"
	"sync"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/scriptpool"
	"github.com/colinyl/lib4go/logger"
)

//RPCScriptProxyHandler 基于脚本的RPC代理服务
type RPCScriptProxyHandler struct {
	tasks          map[string]string
	clusterClient cluster.IClusterClient
	scriptPool    *scriptpool.ScriptPool
	Log           *logger.Logger
	lock          sync.RWMutex
	OnOpenTask    func(task cluster.TaskItem) string
	OnCloseTask   func(task cluster.TaskItem, path string)
}

//NewRPCScriptProxyHandler 构建JOB consumer处理对象
func NewRPCScriptProxyHandler(client cluster.IClusterClient, pool *scriptpool.ScriptPool) *RPCScriptProxyHandler {
	job := &RPCScriptProxyHandler{}
	job.clusterClient = client
	job.scriptPool = pool
	job.tasks = make(map[string]string)
	job.Log, _ = logger.New("job consumer", true)
	return job
}

//GetTasks 获取当前已注册服务列表
func (h *RPCScriptProxyHandler) GetTasks() map[string]string {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.tasks
}

//OpenTask 启动新的任务
func (h *RPCScriptProxyHandler) OpenTask(task cluster.TaskItem) {
	path := task.Name
	if h.OnOpenTask != nil {
		path = h.OnOpenTask(task)
	}
	h.lock.Lock()
	h.tasks[task.Name] = path
	h.lock.Unlock()

}

//CloseTask 关闭任务
func (h *RPCScriptProxyHandler) CloseTask(ti cluster.TaskItem) {
	if path, ok := h.tasks[ti.Name]; ok {
		if h.OnCloseTask != nil {
			h.OnCloseTask(ti, path)
		}
		h.lock.Lock()
		delete(h.tasks, ti.Name)
		h.lock.Unlock()
	}
}

//Request 执行Request请求
func (h *RPCScriptProxyHandler) Request(ti cluster.TaskItem, input string) (string, error) {
	return h.getResult(h.scriptPool.Call(ti.Script, input, ti.Params))
}

//Send 暂不支持
func (h *RPCScriptProxyHandler) Send(ti cluster.TaskItem, input string, data []byte) (string, error) {
	return "", errors.New("job consumer not support send method")
}

//Get 暂不支持
func (h *RPCScriptProxyHandler) Get(ti cluster.TaskItem, input string) ([]byte, error) {
	return nil, errors.New("job consumer not support get method")
}
func (h *RPCScriptProxyHandler) getResult(result []string, er error) (r string, err error) {
	err = er
	if err != nil {
		return
	}
	if len(result) > 0 {
		r = result[0]
	}
	return
}
