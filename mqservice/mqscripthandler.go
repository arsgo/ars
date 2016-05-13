package mqservice

import (
	"strings"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/scriptpool"
)

//MQScriptHandler 脚本处理程序
type MQScriptHandler struct {
	pool *scriptpool.ScriptPool
}

//NewMQScriptHandler 创建新的脚本处理程序
func NewMQScriptHandler(pool *scriptpool.ScriptPool) (mq *MQScriptHandler) {
	return *MQScriptHandler{pool: pool}
}

//Handle 处理MQ消息
func (mq *MQScriptHandler) Handle(task cluster.TaskItem, msg string) bool {
	result, err := mq.pool.Call(task.Name, msg, task.Params)
	if err != nil {
		return false
	}
	return len(result) > 0 && strings.EqualFold(strings.ToLower(result[0]), "true")
}
