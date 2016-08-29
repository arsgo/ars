package base

import (
	"time"

	"github.com/arsgo/lib4go/logger"
)

const (
	TN_MQ_CONSUMER      = "mq.consumer"
	TN_JOB_CONSUMER     = "job.consumer"
	TN_JOB_LOCAL        = "job.local"
	TN_HTTP_API     = "http.api"
	TN_SERVICE_PROVIDER = "service.provider"
)

type InvokeContext struct {
	Session   string
	Input     string
	Params    string
	Body      string
	Log       logger.ILogger
	TaskName  string
	TaskType  string
	StartTime time.Time
}

func NewInvokeContext(taskName string, taskType string, loggerName string, session string, input string, params string, body string) InvokeContext {
	context := InvokeContext{TaskName: taskName, TaskType: taskType, Session: session, Input: input, Params: params, Body: body, StartTime: time.Now()}
	context.Log, _ = logger.NewSession(loggerName, session)
	return context
}
func (c *InvokeContext) PassTime() time.Duration {
	return time.Now().Sub(c.StartTime)
}
