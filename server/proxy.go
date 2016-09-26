package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/arsgo/goproxy"
	"github.com/arsgo/lib4go/logger"
)

//HTTPProxy http代理服务器
type HTTPProxy struct {
	port       string
	Log        logger.ILogger
	loggerName string
}

//NewHTTPProxy 创建http代理服务器
func NewHTTPProxy(port string, lgname string) (h *HTTPProxy) {
	h = &HTTPProxy{loggerName: lgname}
	h.Log, _ = logger.Get(lgname)
	if strings.EqualFold(port, "") {
		h.port = ":8080"
	} else if !strings.HasPrefix(port, ":") {
		h.port = ":" + port
	} else {
		h.port = port
	}
	return h
}

//Start 启动服务器
func (p *HTTPProxy) Start() error {
	proxy := goproxy.NewProxyHttpServer(p.loggerName)
	proxy.Verbose = true
	log.Printf(" -> proxy server: %s...启动完成\n", p.port)
	return http.ListenAndServe(p.port, proxy)
}
