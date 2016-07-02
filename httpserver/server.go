package httpserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/colinyl/ars/cluster"
	"github.com/colinyl/ars/rpcproxy"
	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/webserver"
)

//HttpScriptServer 基于LUA的HTTP服务器
type HttpScriptServer struct {
	Address    string
	handler    []webserver.WebHandler
	routes     []*cluster.ServerRouteConfig
	server     *webserver.WebServer
	Log        logger.ILogger
	loggerName string
	call       func(script string, input string, params string, body string) ([]string, map[string]string, error)
}

//httpScriptController controller
type HttpScriptController struct {
	config *cluster.ServerRouteConfig
	server *HttpScriptServer
}

//NewHttpScriptServer 创建基于LUA的HTTP服务器
func NewHttpScriptServer(Address string, routes []*cluster.ServerRouteConfig, call func(name string, input string, params string, body string) ([]string, map[string]string, error), loggerName string) (server *HttpScriptServer, err error) {
	server = &HttpScriptServer{}
	server.routes = routes
	server.Address = Address
	server.call = call
	server.Log, err = logger.Get(loggerName, true)
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
	r.server = webserver.NewWebServer(r.Address, r.loggerName, r.getHandlers()...)
	go r.server.Serve()
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
func (r *HttpScriptController) getBodyText(request *http.Request) string {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}

//Handle 脚本处理程序(r *HttpScriptController) Handle(ctx *web.Context)
func (r *HttpScriptController) Handle(ctx http.ResponseWriter, request *http.Request) {
	r.server.Log.Info("api.start:", r.config.Script)
	body := r.getBodyText(request)
	request.ParseForm()
	params := make(map[string]string)
	if len(request.Form) > 0 {
		for k, v := range request.Form {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
	}
	data, err := json.Marshal(&params)
	if err != nil {
		r.setResponse(ctx, make(map[string]string), 500, err.Error())
		return
	}

	result, output, err := r.server.call(r.config.Script, string(data), r.config.Params, body)
	r.setHeader(ctx, output)
	if err != nil {
		r.setResponse(ctx, output, 500, err.Error())
		return
	}
	if len(result) == 0 {
		r.setResponse(ctx, output, 200, "")
		return
	}
	if len(result) == 1 {
		r.setResponse(ctx, output, 200, result[0])
		return
	}
	if len(result) == 2 && result[0] == "302" {
		r.setResponse(ctx, output, 302, result[1])
	} else {
		r.setResponse(ctx, output, 500, "system busy")
	}
	return

}
func (r *HttpScriptController) setHeader(ctx http.ResponseWriter, input map[string]string) {
	for i, v := range input {
		ctx.Header().Set(i, v)
	}
}
func (r *HttpScriptController) setResponse(ctx http.ResponseWriter, config map[string]string, code int, msg string) {
	responseContent := ""
	switch code {
	case 200:
		{
			responseContent = rpcproxy.GetDataResult(msg, strings.EqualFold(config["Content-Type"], "text/plain"))
			ctx.Write([]byte(responseContent))
		}
	case 500:
		{
			responseContent = rpcproxy.GetErrorResult(string(code), msg)
			ctx.WriteHeader(500)
			ctx.Write([]byte(responseContent))
		}
	case 302:
		{
			responseContent = msg
			ctx.Header().Set("Location", responseContent)
			ctx.WriteHeader(302)
			ctx.Write([]byte("Redirecting to: " + responseContent))
		}
	}
	r.server.Log.Infof("api.response:%d %s", code, responseContent)
}
