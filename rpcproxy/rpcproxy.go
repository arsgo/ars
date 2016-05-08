package rpcproxy

import (
	"errors"
	"sync"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/rpcclient"
	"github.com/colinyl/lib4go/logger"
)

//ISnap 快照信息接口
type ISnap interface {
	GetSnap() string
}
//RPCClientProxyHandler 处理JOB Consumer操作
type RPCClientProxyHandler struct {
	tasks         map[string]cluster.TaskItem
	clusterClient cluster.IClusterClient
	client        *rpcclient.RPCClient
	Log           *logger.Logger
	snap          ISnap
	lock          sync.RWMutex
}

//NewJobConsumerHandler 构建JOB consumer处理对象
func NewRPCClientProxyHandler(client cluster.IClusterClient, rpcClient *rpcclient.RPCClient, snap ISnap) *RPCClientProxyHandler {
	job := &RPCClientProxyHandler{}
	job.clusterClient = client
	job.client = rpcClient
	job.snap = snap
	job.tasks = make(map[string]cluster.TaskItem)
	job.Log, _ = logger.New("job consumer", true)
	return job
}

//GetTasks 获取当前已注册task列表
func (h *RPCClientProxyHandler) GetTasks() map[string]cluster.TaskItem {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.tasks
}

//OpenTask 启动新的任务
func (h *RPCClientProxyHandler) OpenTask(task cluster.TaskItem) {
	h.lock.Lock()
	h.tasks[task.Name] = task
	h.lock.Unlock()
	h.Log.Info("::start job service:", task.Name)
}

//CloseTask 关闭任务
func (h *RPCClientProxyHandler) CloseTask(ti cluster.TaskItem) {
	if _, ok := h.tasks[ti.Name]; ok {
		h.lock.Lock()
		delete(h.tasks, ti.Name)
		h.lock.Unlock()
	}
}

//Request 执行Request请求
func (h *RPCClientProxyHandler) Request(ti cluster.TaskItem, input string) (string, error) {
	return h.client.Request(ti.Name, input)
}

//Send 暂不支持
func (h *RPCClientProxyHandler) Send(ti cluster.TaskItem, input string, data []byte) (string, error) {
	return "", errors.New("job consumer not support send method")
}

//Get 暂不支持
func (h *RPCClientProxyHandler) Get(ti cluster.TaskItem, input string) ([]byte, error) {
	return nil, errors.New("job consumer not support get method")
}
func (h *RPCClientProxyHandler) getResult(result []string, er error) (r string, err error) {
	err = er
	if err != nil {
		return
	}
	if len(result) > 0 {
		r = result[0]
	}
	return
}
