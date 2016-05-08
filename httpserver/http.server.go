package httpserver

import (
	"encoding/json"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/webserver"
	"github.com/colinyl/web"
)

//HttpScriptServer 基于LUA的HTTP服务器
type HttpScriptServer struct {
	Address string
	handler []webserver.WebHandler
	routes  []*cluster.ServerRouteConfig
	server  *webserver.WebServer
	Log     *logger.Logger
	call    func(script string, input string, params string) ([]string, error)
}

//httpScriptController controller
type httpScriptController struct {
	config cluster.ServerRouteConfig
	server *HttpScriptServer
}

//NewHttpScriptServer 创建基于LUA的HTTP服务器
func NewHttpScriptServer(routes []*cluster.ServerRouteConfig, call func(name string, input string, params string) ([]string, error)) (server *HttpScriptServer, err error) {
	server = &HttpScriptServer{}
	server.routes = routes
	server.call = call
	server.Log, err = logger.New("rpc.server", true)
	return
}

//Stop 停止服务器
func (r *HttpScriptServer) Stop() {
	if r.server != nil {
		r.server.Stop()
	}
}

//Start 启动服务器
func (r *HttpScriptServer) Start() {
	r.Address = rpcservice.GetLocalRandomAddress(20320)
	r.server = webserver.NewWebServer(r.Address, r.getHandlers()...)
	r.server.Serve()
	r.Log.Infof("::start api server%s", r.Address)
}

//getHandlers 获取基于LUA的路由处理程序
func (r *HttpScriptServer) getHandlers() (handlers []webserver.WebHandler) {
	for _, v := range r.routes {
		handler := webserver.WebHandler{Path: v.Path, Method: v.Method}
		handler.Handler = NewHttpScriptController(r).Handle
		handlers = append(handlers, handler)
	}
	return
}

//NewHttpScriptController 创建路由处理程序
func NewHttpScriptController(r *HttpScriptServer) *httpScriptController {
	return &httpScriptController{server: r}
}

//Handle 脚本处理程序
func (r *httpScriptController) Handle(ctx *web.Context) {
	data, err := json.Marshal(&ctx.Params)
	if err != nil {
		r.server.Log.Error(err)
		ctx.Abort(500, err.Error())
		return
	}
	result, err := r.server.call(r.config.Script, string(data), "{}")
	if err != nil {
		r.server.Log.Error(err)
		ctx.Abort(500, err.Error())
		return
	}
	if len(result) == 0 {
		return
	}
	if len(result) == 1 {
		ctx.ResponseWriter.Write([]byte(result[0]))
		return
	}
	if len(result) == 2 && result[0] == "302" {
		ctx.Redirect(302, result[1])
	} else {
		ctx.Redirect(500, "err")
	}
	return

}
