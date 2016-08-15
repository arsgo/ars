package cluster

import (
	"time"

	"github.com/arsgo/lib4go/zkclient"
)

//NewDomainClusterClient  构建集群管理客户端
func NewDomainClusterClient(domain string, localip string, loggerName string, zkips ...string) (c IClusterClient, err error) {
	handler, err := NewClusterHandler(loggerName, zkips...)
	if err != nil {
		return
	}
	c, err = NewClusterClient(domain, localip, handler, loggerName)
	return
}

//NewDomainClusterClientHandler 根据处理程序构建集群管理客户端
func NewDomainClusterClientHandler(domain string, localip string, loggerName string, handler IClusterHandler) (c IClusterClient, err error) {
	c, err = NewClusterClient(domain, localip, handler, loggerName)
	return
}

//NewClusterHandler 构建集群处理程序
func NewClusterHandler(loggerName string, ips ...string) (handler IClusterHandler, err error) {
	handler, err = zkClient.New(ips, time.Second*1, loggerName)
	return
}
