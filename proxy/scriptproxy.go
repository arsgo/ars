package proxy

import (
	"errors"
	"runtime/debug"
	"time"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/ars/script"
	"github.com/arsgo/lib4go/logger"
)

//ScriptProxy 基于脚本的RPC代理服务
type ScriptProxy struct {
	//	tasks         *concurrent.ConcurrentMap
	clusterClient cluster.IClusterClient
	scriptPool    *script.ScriptPool
	Log           logger.ILogger
	loggerName    string
	taskType      string
	OnOpenTask    func(task cluster.TaskItem) string
	OnCloseTask   func(task cluster.TaskItem, path string)
}

//NewScriptProxy 构建JOB consumer处理对象
func NewScriptProxy(client cluster.IClusterClient, pool *script.ScriptPool, taskType string, loggerName string) *ScriptProxy {
	job := &ScriptProxy{loggerName: loggerName, taskType: taskType}
	job.clusterClient = client
	job.scriptPool = pool
	//	job.tasks = concurrent.NewConcurrentMap()
	job.Log, _ = logger.Get(loggerName)
	return job
}

//OpenTask 启动新的任务
func (h *ScriptProxy) OpenTask(task cluster.TaskItem) string {
	path := task.Name
	if h.OnOpenTask != nil {
		path = h.OnOpenTask(task)
	}
	return path
}

//CloseTask 关闭任务
func (h *ScriptProxy) CloseTask(ti cluster.TaskItem, path string) {
	if h.OnCloseTask != nil {
		h.OnCloseTask(ti, path)
	}
}

//Request 执行Request请求
func (h *ScriptProxy) Request(ti cluster.TaskItem, input string, session string, timeout time.Duration) (result string, err error) {
	defer h.recover()
	sresult, smap, err := h.scriptPool.Call(ti.Script, base.NewInvokeContext(ti.Name, h.taskType, h.loggerName, session, input, ti.Params, ""))
	result, _, er := h.getResult(sresult, smap, err)
	if er != nil {
		result = base.GetErrorResult(base.ERR_NOT_FIND_SRVS, er.Error())
	} else {
		result = base.GetDataResult(result, base.IsRaw(smap))
	}
	return
}

//Send 暂不支持
func (h *ScriptProxy) Send(ti cluster.TaskItem, input string, data []byte, timeout time.Duration) (string, error) {
	return "", errors.New("job consumer not support send method")
}

//Get 暂不支持
func (h *ScriptProxy) Get(ti cluster.TaskItem, input string, timeout time.Duration) ([]byte, error) {
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
