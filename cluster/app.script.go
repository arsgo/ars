package cluster

import (
	"errors"
	"fmt"
	"sync"
)

type rpcHandler struct {
	mutex sync.Mutex
	app   *appServer
	queue chan []interface{}
	Name  string
}

func (r *rpcHandler) New() {	
	fmt.Println(r.app)
	fmt.Println("colin")
	r.Name = "colin"
}
func NewRpcHandler(app *appServer) rpcHandler {
	return rpcHandler{queue: make(chan []interface{}, 1), app: app}
}

func (h *rpcHandler) get() (r interface{}, err interface{}) {
	result := <-h.queue
	if len(result) != 2 {
		return "", errors.New("rpc method result value len is error")
	}
	r = result[0]
	err = result[1]
	return
}

func (h *rpcHandler) request(name string, input string) (result string, err error) {
	result, err = h.app.rcServerPool.Request(h.app.rcServicesMap.Next("-"), name, input)
	return
}

func (h *rpcHandler) asyncRequest(name string, input string) {
	h.mutex.Lock()
	if h.queue == nil {
		h.queue = make(chan []interface{}, 1)
	}
	h.mutex.Unlock()
	go func(h *rpcHandler, name string, input string) {
		result, err := h.request(name, input)
		if err != nil {
			h.queue <- []interface{}{result, err.Error()}
		} else {
			h.queue <- []interface{}{result, nil}
		}

	}(h, name, input)
	return
}
