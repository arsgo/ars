package httpserver

import (
	"encoding/json"
	"strings"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/rpcproxy"
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
	call    func(script string, input string, params string) ([]string,map[string]string,error)
}

//httpScriptController controller
type HttpScriptController struct {
	config *cluster.ServerRouteConfig
	server *HttpScriptServer
}

//NewHttpScriptServer 创建基于LUA的HTTP服务器
func NewHttpScriptServer(Address string, routes []*cluster.ServerRouteConfig, call func(name string, input string, params string) ([]string,map[string]string, error)) (server *HttpScriptServer, err error) {
	server = &HttpScriptServer{}
	server.routes = routes
	server.Address = Address
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
	if strings.EqualFold(r.Address, "") {
		r.Address = rpcservice.GetLocalRandomAddress(20320)
	} else if !strings.HasPrefix(r.Address, ":") {
		r.Address = ":" + r.Address
	}
	r.server = webserver.NewWebServer(r.Address, r.getHandlers()...)
	r.server.Serve()
	r.Log.Infof("::start api server%s", r.Address)
}

//getHandlers 获取基于LUA的路由处理程序
func (r *HttpScriptServer) getHandlers() (handlers []webserver.WebHandler) {
	for _, v := range r.routes {
		handler := webserver.WebHandler{Path: v.Path, Method: v.Method, Script: v.Script}
		handler.Handler = NewHttpScriptController(r, v).Handle
		handlers = append(handlers, handler)
	}
	return
}

//NewHttpScriptController 创建路由处理程序
func NewHttpScriptController(r *HttpScriptServer, config *cluster.ServerRouteConfig) *HttpScriptController {
	return &HttpScriptController{server: r, config: config}
}

//Handle 脚本处理程序
func (r *HttpScriptController) Handle(ctx *web.Context) {
	r.server.Log.Info("api.start:", r.config.Script)
	data, err := json.Marshal(&ctx.Params)
	if err != nil {
		r.setResponse(ctx, 500, err.Error())
		return
	}
	result,output, err := r.server.call(r.config.Script, string(data), r.config.Params)
	r.setHeader(ctx,output)
	if err != nil {
		r.setResponse(ctx, 500, err.Error())
		return
	}
	if len(result) == 0 {
		r.setResponse(ctx, 200, "")
		return
	}
	if len(result) == 1 {
		r.setResponse(ctx, 200, result[0])
		return
	}
	if len(result) == 2 && result[0] == "302" {
		r.setResponse(ctx, 302, result[1])
	} else {
		r.setResponse(ctx, 500, "system busy")
	}
	return

}
func (r *HttpScriptController) setHeader(ctx *web.Context,input map[string]string){
	for i,v:=range input{
		ctx.SetHeader(i,v,true)
	}
}
func (r *HttpScriptController) setResponse(ctx *web.Context, code int, msg string) {

	responseContent := ""
	switch code {
	case 200:
		{
			responseContent = rpcproxy.GetDataResult(msg)
			ctx.ResponseWriter.Write([]byte(responseContent))
		}
	case 500:
		{
			responseContent = rpcproxy.GetErrorResult(string(code), msg)
			ctx.Abort(500, responseContent)
		}
	case 302:
		{
			responseContent = msg
			ctx.Redirect(302, responseContent)
		}
	}
	r.server.Log.Infof("api.response:%d %s", code, responseContent)
}
