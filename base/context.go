package base

import "github.com/colinyl/lib4go/logger"

type InvokeContext struct {
	Session string
	Input   string
	Params  string
	Body    string
	Log     logger.ILogger
}

func NewInvokeContext(session string, input string, params string, body string) InvokeContext {
	context := InvokeContext{Session: session, Input: input, Params: params, Body: body}
	context.Log = logger.GetDeubgLogger(session)
	return context

}
