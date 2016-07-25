package proxy

import (
	"errors"
	"runtime/debug"

	"github.com/colinyl/ars/base"
	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/script"
	"github.com/colinyl/lib4go/concurrent"
	"github.com/colinyl/lib4go/logger"
)

//ScriptProxy 基于脚本的RPC代理服务
type ScriptProxy struct {
	tasks         concurrent.ConcurrentMap
	clusterClient cluster.IClusterClient
	scriptPool    *script.ScriptPool
	Log           logger.ILogger
	OnOpenTask    func(task cluster.TaskItem) string
	OnCloseTask   func(task cluster.TaskItem, path string)
}

//NewScriptProxy 构建JOB consumer处理对象
func NewScriptProxy(client cluster.IClusterClient, pool *script.ScriptPool, loggerName string) *ScriptProxy {
	job := &ScriptProxy{}
	job.clusterClient = client
	job.scriptPool = pool
	job.tasks = concurrent.NewConcurrentMap()
	job.Log, _ = logger.Get(loggerName)
	return job
}

//GetTasks 获取当前已注册服务列表
func (h *ScriptProxy) GetTasks() map[string]string {
	data := make(map[string]string)
	service := h.tasks.GetAll()
	for i, v := range service {
		data[i] = v.(string)
	}
	return data
}

//OpenTask 启动新的任务
func (h *ScriptProxy) OpenTask(task cluster.TaskItem) {
	path := task.Name
	if h.OnOpenTask != nil {
		path = h.OnOpenTask(task)
	}
	h.tasks.Set(task.Name, path)
}

//CloseTask 关闭任务
func (h *ScriptProxy) CloseTask(ti cluster.TaskItem) {
	value := h.tasks.Get(ti.Name)
	if value != nil && h.OnCloseTask != nil {
		h.OnCloseTask(ti, value.(string))
	}
	h.tasks.Delete(ti.Name)
}

//Request 执行Request请求
func (h *ScriptProxy) Request(ti cluster.TaskItem, input string, session string) (result string, err error) {
	defer h.recover()
	sresult, smap, err := h.scriptPool.Call(ti.Script, base.NewInvokeContext(session, input, ti.Params, ""))
	result, _, er := h.getResult(sresult, smap, err)
	if er != nil {
		result = base.GetErrorResult("500", er.Error())
	} else {
		result = base.GetDataResult(result, base.IsRaw(smap))
	}
	return
}

//Send 暂不支持
func (h *ScriptProxy) Send(ti cluster.TaskItem, input string, data []byte) (string, error) {
	return "", errors.New("job consumer not support send method")
}

//Get 暂不支持
func (h *ScriptProxy) Get(ti cluster.TaskItem, input string) ([]byte, error) {
	return nil, errors.New("job consumer not support get method")
}
func (h *ScriptProxy) getResult(result []string, params map[string]string, er error) (r string, p map[string]string, err error) {
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
func (h *ScriptProxy) recover() {
	if r := recover(); r != nil {
		h.Log.Fatal(r, string(debug.Stack()))
	}
}
