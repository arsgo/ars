package webservice

import (
	"github.com/colinyl/web"
)

//WebHandler Web处理程序
type WebHandler struct {
	Path    string
	Method  string
	Handler interface{}
}

//WebService WEB服务
type WebService struct {
	routes  []WebHandler
	server  *web.Server
	address string
}

//NewWebService 创建WebService服务
func NewWebService(address string, handlers ...WebHandler) (service *WebService) {
	service = &WebService{}
	service.server = web.NewServer()
	service.routes = handlers
	service.address = address
	return
}

//Serve 启动WEB服务器
func (w *WebService) Serve() {
	for _, handler := range w.routes {
		switch handler.Method {
		case "get":
			w.server.Get(handler.Path, handler.Handler)
		case "post":
			w.server.Post(handler.Path, handler.Handler)
		case "put":
			w.server.Put(handler.Path, handler.Handler)
		case "delete":
			w.server.Delete(handler.Path, handler.Handler)
		case "*":
			w.server.Get(handler.Path, handler.Handler)
			w.server.Post(handler.Path, handler.Handler)
		}
	}
	go w.server.Run(w.address)
}

//Stop 停止服务器
func (w *WebService) Stop() {	
	w.server.Close()
}
