package cluster

import l "github.com/yuin/gopher-lua"

func (a *appServer) asyncRequest(L *l.LState) (handler *scriptCallbackHandler) {
	name := L.ToString(1)
	input := L.ToString(2)
	handler = &scriptCallbackHandler{queue: make(chan []string, 1), Log: a.Log}
	go func(name string, input string) {
		handler.queue <- a.request(name, input)
	}(name, input)
	return
}
func (a *appServer) asyncSend(L *l.LState) (handler *scriptCallbackHandler) {
	handler = &scriptCallbackHandler{queue: make(chan []string, 1), Log: a.Log}
	go func() {
		handler.queue <- a.send(L)
	}()
	return
}

func (a *appServer) request(name string, input string) (result []string) {
	rest, err := a.rcServerPool.Request(a.rcServicesMap.Next("-"), name, input)
	result = append(result, rest)
	if err != nil {
		result = append(result, err.Error())
	} else {
		result = append(result, "")
	}
	return
}

func (a *appServer) send(L *l.LState) (result []string) {
	name := L.ToString(1)
	input := L.ToString(2)
	buffer := []byte(L.ToString(3))
	group := a.rcServicesMap.Next("-")
	rest, err := a.rcServerPool.Send(group, name, input, buffer)
	result = append(result, rest)
	if err != nil {
		result = append(result, err.Error())
	} else {
		result = append(result, "")
	}
	return
}
