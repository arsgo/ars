package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/arsgo/ars/base"
	"github.com/arsgo/ars/base/rpcservice"
	"github.com/arsgo/ars/cluster"
	"github.com/arsgo/lib4go/logger"
	"github.com/arsgo/lib4go/webserver"
)

//HTTPScriptServer 基于LUA的HTTP服务器
type HTTPScriptServer struct {
	Address    string
	handler    []webserver.WebHandler
	routes     []*cluster.ServerRouteConfig
	server     *webserver.WebServer
	Log        logger.ILogger
	snap       *base.ServerSnap
	loggerName string
	call       func(string, base.InvokeContext) ([]string, map[string]string, error)
}

//HTTPScriptController controller
type HTTPScriptController struct {
	config     *cluster.ServerRouteConfig
	server     *HTTPScriptServer
	snap       *base.ServerSnap
	loggerName string
}

//NewHTTPScriptServer 创建基于LUA的HTTP服务器
func NewHTTPScriptServer(Address string, routes []*cluster.ServerRouteConfig, call func(string, base.InvokeContext) ([]string, map[string]string, error), loggerName string) (server *HTTPScriptServer, err error) {
	server = &HTTPScriptServer{snap: &base.ServerSnap{}}
	server.routes = routes
	server.Address = Address
	server.call = call
	server.Log, err = logger.Get(loggerName)
	return
}

//Stop 停止服务器
func (r *HTTPScriptServer) Stop() {
	if r.server != nil {
		r.server.Stop()
	}
}

//Start 启动服务器
func (r *HTTPScriptServer) Start() {
	if strings.EqualFold(r.Address, "") {
		r.Address = rpcservice.GetLocalRandomAddress(20320)
	} else if !strings.HasPrefix(r.Address, ":") {
		r.Address = ":" + r.Address
	}
	r.server = webserver.NewWebServer(r.Address, r.loggerName, r.getHandlers()...)
	go func() {
		er := r.server.Serve()
		r.Log.Error(er)
	}()

	r.Log.Infof("::start http server:%s", r.Address)
}

//GetSnap 获取当前服务器快照信息
func (r *HTTPScriptServer) GetSnap() base.ServerSnap {
	return *r.snap
}

//getHandlers 获取基于LUA的路由处理程序
func (r *HTTPScriptServer) getHandlers() (handlers []webserver.WebHandler) {
	for _, v := range r.routes {
		handler := webserver.WebHandler{Path: v.Path, Encoding: v.Encoding, Method: v.Method, Script: v.Script, LoggerName: r.loggerName}
		handler.Handler = NewHTTPScriptController(r, v, r.snap).Handle
		handlers = append(handlers, handler)
	}
	return
}

//NewHTTPScriptController 创建路由处理程序
func NewHTTPScriptController(r *HTTPScriptServer, config *cluster.ServerRouteConfig, snap *base.ServerSnap) *HTTPScriptController {
	return &HTTPScriptController{server: r, config: config, snap: snap, loggerName: r.loggerName}
}

func (r *HTTPScriptController) getBodyText(encoding string, request *http.Request) (content string, err error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}
	content = changeEncodingData(encoding, body)
	return
}
func (r *HTTPScriptController) getPostValues(encoding string, body string) (rt map[string]string, err error) {
	rt = make(map[string]string)
	values, err := url.ParseQuery(body)
	if err != nil {
		return
	}
	for i, v := range values {
		if len(v) > 0 && !strings.EqualFold(v[0], "") {
			rt[i] = changeEncodingData(encoding, []byte(v[0]))
		}
	}
	return
}

//Handle 脚本处理程序(r *HttpScriptController) Handle(ctx *web.Context)
func (r *HTTPScriptController) Handle(context *webserver.Context) {
	defer r.snap.Add(time.Now())
	body, err := r.getBodyText(context.Encoding, context.Request)
	if err != nil {
		context.Log.Error("获取POST数据异常:", err)
	}
	context.Request.ParseForm()
	params, err := r.getPostValues(context.Encoding, body)
	if err != nil {
		context.Log.Error("解析POST数据异常:", err)
	}
	if len(context.Request.Form) > 0 {
		for k, v := range context.Request.Form {
			if len(v) > 0 && len(v[0]) > 0 && !strings.EqualFold(v[0], "") {
				params[k] = changeEncodingData(context.Encoding, []byte(v[0]))
			}
		}
	}
	data, err := json.Marshal(&params)
	if err != nil {
		context.Log.Info("-->api.request/response.error:", context.Request.URL.Path, err)
		r.setResponse(context, make(map[string]string), 500, err.Error())
		return
	}
	input := string(data)
	context.Log.Info("-->api.request:", context.Request.URL.Path, input, body)
	result, output, err := r.server.call(r.config.Script, base.NewInvokeContext(r.loggerName, context.Session, input, r.config.Params, body))
	r.setHeader(context.Writer, output)
	if err != nil {
		r.setResponse(context, output, 500, err.Error())
		return
	}
	switch len(result) {
	case 0:
		r.setResponse(context, output, 200, "")
	case 1:
		r.setResponse(context, output, 200, result[0])
	case 2:
		if result[0] == "302" {
			r.setResponse(context, output, 302, result[1])
		} else {
			r.setResponse(context, output, 500, "system busy")
		}
	default:
		r.setResponse(context, output, 500, "system busy")
	}

	return

}
func (r *HTTPScriptController) setHeader(ctx http.ResponseWriter, input map[string]string) {
	for i, v := range input {
		if strings.HasPrefix(i, "_") {
			continue
		}
		ctx.Header().Set(i, v)
	}
}
func (r *HTTPScriptController) setResponse(context *webserver.Context, config map[string]string, code int, msg string) {
	responseContent := ""
	switch code {
	case 200:
		{
			responseContent = base.GetDataResult(msg, base.IsRaw(config))
			context.Writer.Write([]byte(responseContent))
		}
	case 500:
		{
			responseContent = base.GetErrorResult(string(code), msg)
			context.Writer.WriteHeader(500)
			context.Writer.Write([]byte(responseContent))
		}
	case 302:
		{
			responseContent = msg
			context.Writer.Header().Set("Location", responseContent)
			context.Writer.WriteHeader(302)
			context.Writer.Write([]byte("Redirecting to: " + responseContent))
		}
	}
	context.Log.Infof("api.response:[%d,%v]%s", code, context.PassTime(), responseContent)
}
func changeEncodingData(encoding string, data []byte) (content string) {
	if !strings.EqualFold(encoding, "GBK") && !strings.EqualFold(encoding, "GB2312") {
		content = string(data)
		return
	}
	buffer, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		content = string(data)
		return
	}
	content = string(buffer)
	return
}
