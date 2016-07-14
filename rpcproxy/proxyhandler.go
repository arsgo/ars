package rpcproxy

import (
	"sync"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/lib4go/concurrent"
	"github.com/colinyl/lib4go/logger"
)

//ISnap 快照信息接口
type ISnap interface {
	GetSnap() string
}

//RPCProxyHandler 处理JOB Consumer操作
type RPCProxyHandler struct {
	tasks         concurrent.ConcurrentMap
	clusterClient cluster.IClusterClient
	client        *RPCClient
	Log           logger.ILogger
	snap          ISnap
	lock          sync.RWMutex
}

//NewRPCProxyHandler 构建JOB consumer处理对象
func NewRPCProxyHandler(client cluster.IClusterClient, rpcClient *RPCClient, snap ISnap, loggerName string) *RPCProxyHandler {
	job := &RPCProxyHandler{}
	job.clusterClient = client
	job.client = rpcClient
	job.snap = snap
	job.tasks = concurrent.NewConcurrentMap()
	job.Log, _ = logger.Get(loggerName, true)
	return job
}

//GetTasks 获取当前已注册task列表
func (h *RPCProxyHandler) GetTasks() map[string]cluster.TaskItem {
	data := make(map[string]cluster.TaskItem)
	service := h.tasks.GetAll()
	for i, v := range service {
		data[i] = v.(cluster.TaskItem)
	}
	return data
}

//OpenTask 启动新的任务
func (h *RPCProxyHandler) OpenTask(task cluster.TaskItem) {
	h.tasks.Set(task.Name, task)
	h.Log.Info("::start service:", task.Name)
}

//CloseTask 关闭任务
func (h *RPCProxyHandler) CloseTask(ti cluster.TaskItem) {
	h.Log.Info(" -> close service:", ti.Name)
	h.tasks.Delete(ti.Name)
}

//Request 执行Request请求
func (h *RPCProxyHandler) Request(ti cluster.TaskItem, input string,session string) (r string, err error) {
	defer h.recover()
	r, _ = h.client.Request(ti.Name, input,session)
	return
}

//Send 暂不支持
func (h *RPCProxyHandler) Send(ti cluster.TaskItem, input string, data []byte) (string, error) {
	return h.client.Send(ti.Name, input, string(data))
}

//Get 暂不支持
func (h *RPCProxyHandler) Get(ti cluster.TaskItem, input string) ([]byte, error) {
	data, err := h.client.Get(ti.Name, input)
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}
func (h *RPCProxyHandler) getResult(result []string, er error) (r string, err error) {
	err = er
	if err != nil {
		return
	}
	if len(result) > 0 {
		r = result[0]
	}
	return
}

func (h *RPCProxyHandler) recover() {
	if r := recover(); r != nil {
		h.Log.Fatal(r)
	}
}
