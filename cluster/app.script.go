package cluster

import (
	"errors"
	"fmt"
	"sync"

	"github.com/colinyl/lib4go/utility"
)

type rpcHandler struct {
	app    *appServer
	queues map[string]chan []interface{}
	mutex  sync.Mutex
}

func (a *appServer) NewRpcHandler() (h *rpcHandler) {
	return &rpcHandler{queues: make(map[string]chan []interface{}), app: a}
}

func (h *rpcHandler) GetAsyncResult(session string) (r interface{}, err interface{}) {
	h.mutex.Lock()
	if _, ok := h.queues[session]; !ok {
		err = errors.New(fmt.Sprint("not find session:", session))
		return
	}
	queue := h.queues[session]
	h.mutex.Unlock()
	result := <-queue
	if len(result) != 2 {
		return "", errors.New("rpc method result value len is error")
	}
	r = result[0]
	err = result[1]
	return
}

func (h *rpcHandler) Request(name string, input string) (result string, err error) {
	result, err = h.app.rcServerPool.Request(h.app.rcServicesMap.Next("-"), name, input)
	return
}
func (h *rpcHandler) Send(name string, input string, data string) (result string, err error) {
	result, err = h.app.rcServerPool.Send(h.app.rcServicesMap.Next("-"), name, input, []byte(data))
	return
}
func (h *rpcHandler) Get(name string, input string) (result string, err error) {
	data, err := h.app.rcServerPool.Get(h.app.rcServicesMap.Next("-"), name, input)
	if err != nil {
		result = string(data)
	}
	return
}

func (h *rpcHandler) AsyncRequest(name string, input string) (session string) {
	session = utility.GetGUID()
	h.mutex.Lock()
	h.queues[session] = make(chan []interface{}, 1)
	h.mutex.Unlock()
	go func(session string, h *rpcHandler, name string, input string) {
		result, err := h.Request(name, input)
		h.mutex.Lock()
		defer h.mutex.Unlock()
		if err != nil {

			h.queues[session] <- []interface{}{result, err.Error()}
		} else {
			h.queues[session] <- []interface{}{result, nil}
		}

	}(session, h, name, input)
	return
}
func (h *rpcHandler) AsyncSend(name string, input string, data string) (session string) {
	session = utility.GetGUID()
	h.mutex.Lock()
	h.queues[session] = make(chan []interface{}, 1)
	h.mutex.Unlock()
	go func(session string, h *rpcHandler, name string, input string, data string) {
		result, err := h.Send(name, input, data)
		h.mutex.Lock()
		defer h.mutex.Unlock()
		if err != nil {
			h.queues[session] <- []interface{}{result, err.Error()}
		} else {
			h.queues[session] <- []interface{}{result, nil}
		}

	}(session, h, name, input, data)
	return
}
func (h *rpcHandler) AsyncGet(name string, input string) (session string) {
	session = utility.GetGUID()
	h.mutex.Lock()
	h.queues[session] = make(chan []interface{}, 1)
	h.mutex.Unlock()
	go func(session string, h *rpcHandler, name string, input string) {
		h.mutex.Lock()
		defer h.mutex.Unlock()
		result, err := h.Get(name, input)
		if err != nil {
			h.queues[session] <- []interface{}{result, err.Error()}
		} else {
			h.queues[session] <- []interface{}{result, nil}
		}

	}(session, h, name, input)
	return
}
