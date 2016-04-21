package cluster

import (
	"encoding/json"

	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/ars/webservice"
	"github.com/colinyl/web"
)

func (r *appServer) StopAPIServer() {
	if r.apiServer != nil {
		r.apiServer.Stop()
	}
}
func (r *appServer) StartAPIServer() {
	 r.scriptHandlers=nil
	 r.scriptHandlers=make(map[string]*scriptHandler)
	r.apiServerAddress = rpcservice.GetLocalRandomAddress()
	r.apiServer = webservice.NewWebService(r.apiServerAddress, r.getAPIServerHandler()...)
	r.apiServer.Serve()
	r.Log.Infof("::start api server%s", r.apiServerAddress)
}
func (r *appServer) getAPIServerHandler() (handlers []webservice.WebHandler) {
	for _, v := range r.scriptServer {
		handler := &scriptHandler{data: v, server: r}
		r.scriptHandlers[v.Path] = handler
		handlers = append(handlers,webservice.WebHandler{v.Path, v.Method, handler.ExecuteScript})
	}
	return
}

func (r *scriptHandler) ExecuteScript(ctx *web.Context) {
	r.server.Log.Infof(">api execute script:%s", r.data.Script)
	data, err := json.Marshal(&ctx.Params)
	r.server.Log.Info(string(data))
	if err != nil {
		r.server.Log.Error(err)
		ctx.Abort(500, err.Error())
		return
	}
	result, err := r.server.scriptEngine.Call(r.data.Script, string(data))
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
