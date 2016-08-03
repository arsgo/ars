package cluster

import (
	"time"

	"github.com/arsgo/lib4go/zkclient"
)

func GetClusterClient(domain string, localip string, loggerName string, zkips ...string) (c IClusterClient, err error) {
	handler, err := zkClient.New(zkips, time.Second*1, loggerName)
	if err != nil {
		return
	}
	c, err = NewClusterClient(domain, localip, handler, loggerName)
	return
}
