package mq

import (
	"strings"

	"github.com/colinyl/ars/base"
	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/script"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

//MQScriptHandler 脚本处理程序
type MQScriptHandler struct {
	pool       *script.ScriptPool
	Log        logger.ILogger
	loggerName string
}

//NewMQScriptHandler 创建新的脚本处理程序
func NewMQScriptHandler(pool *script.ScriptPool, loggerName string) (mq *MQScriptHandler) {
	mq = &MQScriptHandler{pool: pool, loggerName: loggerName}
	mq.Log, _ = logger.Get(loggerName)
	return
}

//Handle 处理MQ消息
func (mq *MQScriptHandler) Handle(task cluster.TaskItem, input string, session string) bool {
	mq.Log.Infof(" -> recv mq message:%s", input)
	result, _, err := mq.pool.Call(task.Script, base.NewInvokeContext(mq.loggerName, utility.GetSessionID(), input, task.Params, ""))
	if err != nil {
		return false
	}
	return len(result) > 0 && (strings.EqualFold(strings.ToLower(result[0]), "true") || strings.EqualFold(strings.ToLower(result[0]), "success"))
}
