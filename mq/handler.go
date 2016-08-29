package mq

import (
	"strings"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/ars/script"
	"github.com/arsgo/lib4go/logger"
	"github.com/arsgo/lib4go/utility"
)

//MQScriptHandler 脚本处理程序
type MQScriptHandler struct {
	pool        *script.ScriptPool
	Log         logger.ILogger
	loggerName  string
	onOpenTask  func(task cluster.TaskItem) string
	onCloseTask func(task cluster.TaskItem, path string)
	collector   base.ICollector
}

//NewMQScriptHandler 创建新的脚本处理程序
func NewMQScriptHandler(pool *script.ScriptPool, loggerName string, onOpenTask func(task cluster.TaskItem) string, onCloseTask func(task cluster.TaskItem, path string), collector base.ICollector) (mq *MQScriptHandler) {
	mq = &MQScriptHandler{pool: pool, loggerName: loggerName}
	mq.onOpenTask = onOpenTask
	mq.onCloseTask = onCloseTask
	mq.collector = collector
	mq.Log, _ = logger.Get(loggerName)
	return
}
func (mq *MQScriptHandler) OnOpenTask(task cluster.TaskItem) string {
	return mq.onOpenTask(task)
}
func (mq *MQScriptHandler) OnCloseTask(task cluster.TaskItem, path string) {
	mq.onCloseTask(task, path)
}

//Handle 处理MQ消息
func (mq *MQScriptHandler) Handle(task cluster.TaskItem, input string, session string) bool {
	context := base.NewInvokeContext(task.Name, base.TN_MQ_CONSUMER, mq.loggerName, utility.GetSessionID(), input, task.Params, "")
	context.Log.Infof("--> mq.request(%s):%s,%s", task.Name, task.Script, input)
	result, _, err := mq.pool.Call(task.Script, context)
	defer context.Log.Infof("--> mq.response(%s,%v):%v", task.Name, context.PassTime(), result)
	if err != nil {
		mq.collector.Error(task.Name)
		return false
	}
	v := len(result) > 0 && (strings.EqualFold(result[0], "true") || strings.EqualFold(result[0], "success"))
	mq.collector.Juge(v, task.Name)
	return v
}
