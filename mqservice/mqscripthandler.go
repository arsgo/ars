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
	Log  *logger.Logger
}

//NewMQScriptHandler 创建新的脚本处理程序
func NewMQScriptHandler(pool *rpcproxy.ScriptPool) (mq *MQScriptHandler) {
	mq = &MQScriptHandler{pool: pool}
	mq.Log, _ = logger.New("mq consumer handler", true)
	return
}

//Handle 处理MQ消息
func (mq *MQScriptHandler) Handle(task cluster.TaskItem, input string) bool {
	mq.Log.Info(" -> recv mq message:", input)
	result, _, err := mq.pool.Call(task.Name, input, task.Params)
	if err != nil {
		return false
	}
	return len(result) > 0 && strings.EqualFold(strings.ToLower(result[0]), "true")
}
