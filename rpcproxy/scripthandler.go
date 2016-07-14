package rpcproxy

import (
	"errors"

	"github.com/colinyl/ars/base"
	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/concurrent"
	"github.com/colinyl/lib4go/logger"
)

//RPCScriptHandler 基于脚本的RPC代理服务
type RPCScriptHandler struct {
	tasks         concurrent.ConcurrentMap
	clusterClient cluster.IClusterClient
	scriptPool    *ScriptPool
	Log           logger.ILogger
	OnOpenTask    func(task cluster.TaskItem) string
	OnCloseTask   func(task cluster.TaskItem, path string)
}

//NewRPCScriptHandler 构建JOB consumer处理对象
func NewRPCScriptHandler(client cluster.IClusterClient, pool *ScriptPool, loggerName string) *RPCScriptHandler {
	job := &RPCScriptHandler{}
	job.clusterClient = client
	job.scriptPool = pool
	job.tasks = concurrent.NewConcurrentMap()
	job.Log, _ = logger.Get(loggerName, true)
	return job
}

//GetTasks 获取当前已注册服务列表
func (h *RPCScriptHandler) GetTasks() map[string]string {
	data := make(map[string]string)
	service := h.tasks.GetAll()
	for i, v := range service {
		data[i] = v.(string)
	}
	return data
}

//OpenTask 启动新的任务
func (h *RPCScriptHandler) OpenTask(task cluster.TaskItem) {
	path := task.Name
	if h.OnOpenTask != nil {
		path = h.OnOpenTask(task)
	}
	h.tasks.Set(task.Name, path)
}

//CloseTask 关闭任务
func (h *RPCScriptHandler) CloseTask(ti cluster.TaskItem) {
	value := h.tasks.Get(ti.Name)
	if value != nil && h.OnCloseTask != nil {
		h.OnCloseTask(ti, value.(string))
	}
	h.tasks.Delete(ti.Name)
}

//Request 执行Request请求
func (h *RPCScriptHandler) Request(ti cluster.TaskItem, input string, session string) (result string, err error) {
	defer h.recover()
	sresult, smap, err := h.scriptPool.Call(ti.Script, base.NewInvokeContext(session, input, ti.Params, ""))
	result, _, er := h.getResult(sresult, smap, err)
	if er != nil {
		result = GetErrorResult("500", er.Error())
	} else {
		result = GetDataResult(result, IsRaw(smap))
	}
	return
}

//Send 暂不支持
func (h *RPCScriptHandler) Send(ti cluster.TaskItem, input string, data []byte) (string, error) {
	return "", errors.New("job consumer not support send method")
}

//Get 暂不支持
func (h *RPCScriptHandler) Get(ti cluster.TaskItem, input string) ([]byte, error) {
	return nil, errors.New("job consumer not support get method")
}
func (h *RPCScriptHandler) getResult(result []string, params map[string]string, er error) (r string, p map[string]string, err error) {
	err = er
	if err != nil {
		return
	}
	if len(result) > 0 {
		r = result[0]
	}
	p = params
	return
}
func (h *RPCScriptHandler) recover() {
	if r := recover(); r != nil {
		h.Log.Fatal(r)
	}
}
