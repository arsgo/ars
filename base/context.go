package base

import (
	"time"

	"github.com/arsgo/lib4go/logger"
)

type InvokeContext struct {
	Session   string
	Input     string
	Params    string
	Body      string
	Log       logger.ILogger
	StartTime time.Time
}

func NewInvokeContext(loggerName string, session string, input string, params string, body string) InvokeContext {
	context := InvokeContext{Session: session, Input: input, Params: params, Body: body, StartTime: time.Now()}
	context.Log, _ = logger.NewSession(loggerName, session)
	return context
}
func (c *InvokeContext) PassTime() time.Duration {
	return time.Now().Sub(c.StartTime)
}
