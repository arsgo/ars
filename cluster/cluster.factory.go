package cluster

import (
	"time"

	"github.com/colinyl/lib4go/zkclient"
)

func GetClusterClient(domain string, localip string, ips ...string) (c IClusterClient, err error) {
	handler, err := zkClient.New(ips, time.Second*1)
	if err != nil {
		return
	}
	c, err = NewClusterClient(domain, localip, handler)
	return
}
