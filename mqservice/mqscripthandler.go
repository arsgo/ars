package mqservice

import (
	"strings"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/rpcproxy"
	"github.com/colinyl/lib4go/logger"
)

//MQScriptHandler 脚本处理程序
type MQScriptHandler struct {
	pool *rpcproxy.ScriptPool
	Log  logger.ILogger
}

//NewMQScriptHandler 创建新的脚本处理程序
func NewMQScriptHandler(pool *rpcproxy.ScriptPool, loggerName string) (mq *MQScriptHandler) {
	mq = &MQScriptHandler{pool: pool}
	mq.Log, _ = logger.Get(loggerName, true)
	return
}

//Handle 处理MQ消息
func (mq *MQScriptHandler) Handle(task cluster.TaskItem, input string) bool {
	mq.Log.Infof(" -> recv mq message:%s", input)
	result, _, err := mq.pool.Call(task.Script, input, task.Params,"")
	if err != nil {
		return false
	}
	return len(result) > 0 && (strings.EqualFold(strings.ToLower(result[0]), "true") || strings.EqualFold(strings.ToLower(result[0]), "success"))
}
