package cluster

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/arsgo/lib4go/utility"
)

var servers string
var clusterClient *ClusterClient
var clusterHandler IClusterHandler
var domain string
var domainName string

func init() {
	servers = "192.168.101.166:2181"
	domain = "/ars/test"
	domainName = "@ars.test"
	clusterHandler, err := NewClusterHandler("cluster.test", servers)
	if err != nil {
		fmt.Println(err)
		os.Exit(100)
	}
	clusterClient, err = NewClusterClient(domain, utility.GetLocalIPAddress("192;172"), clusterHandler, "cluster.test")
	if err != nil {
		fmt.Println(err)
		os.Exit(100)
	}
}

func TestDomainValues(t *testing.T) {
	domain := clusterClient.getDomain("ars/test/")
	expect := "/ars/test"
	if domain != expect {
		t.Errorf("domain error:actual:%s,expect:%s", domain, expect)
	}
	domain = clusterClient.getDomain("/ars/test/")
	if domain != expect {
		t.Errorf("domain error:actual:%s,expect:%s", domain, expect)
	}
	domain = clusterClient.getDomain("/ars/test")
	if domain != expect {
		t.Errorf("domain error:actual:%s,expect:%s", domain, expect)
	}
}
func TestDomainNameValues(t *testing.T) {
	domainName := clusterClient.getDomainName("ars/test/")
	expect := "@ars.test"
	if domainName != expect {
		t.Errorf("domain error:actual:%s,expect:%s", domainName, expect)
	}
	domainName = clusterClient.getDomainName("/ars/test/")
	if domainName != expect {
		t.Errorf("domain error:actual:%s,expect:%s", domainName, expect)
	}
	domainName = clusterClient.getDomainName("ars/test.a")
	expect = "@ars.test.a"
	if domainName != expect {
		t.Errorf("domain error:actual:%s,expect:%s", domainName, expect)
	}
}
func TestWatchPathExists(t *testing.T) {
	t.Log("监控：节点不存在是回调通知")
	existsChan := make(chan bool, 1)
	path := "/ars/test/app/1245658745541"
	go clusterClient.WaitClusterPathExists(path, time.Hour, func(path string, exists bool) {
		existsChan <- exists
	})
	select {
	case v := <-existsChan:
		if v {
			t.Error(path, "节点不存在，应该立即返回false,但实际值为:", v)
		}
	}
	t.Log("监控：节点已存在时回调通知")
	path, err := clusterClient.handler.CreateTmpNode(path, "")
	if err != nil {
		t.Error("创建临时节点失败:", err)
	}
	select {
	case <-time.After(time.Second * 3):
		t.Error("创建创建成功WaitClusterPathExists方法应在3秒内返回节点已存在，但未返回")
	case v := <-existsChan:
		if !v {
			t.Error(path, "创建创建成功后应该返回true,但实际返回了false")
		}
	}

	path2 := "/ars/test/app"
	go clusterClient.WaitClusterPathExists(path2, time.Hour, func(path string, exists bool) {
		existsChan <- exists
	})
	select {
	case v := <-existsChan:
		if !v {
			t.Error(path, "节点已存在，应该立即返回true,但实际值为:", v)
		}
	}
}
func TestWatchValueChange(t *testing.T) {
	path := "/cluster/test/app/test"
	value := `{"id":1}`
	path, err := clusterClient.handler.CreateTmpNode(path, value)
	if err != nil {
		t.Error("创建临时节点失败:", err)
	}
}
