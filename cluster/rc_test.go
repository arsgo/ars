package cluster

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func CheckServiceChange(server *rcServer, t *testing.T) {
	time.Sleep(time.Second * 5)
}

func TMasterServer(server *rcServer, pt int64, spCount int, t *testing.T) {
	if !server.IsMasterServer {
		t.Error("master bind failed")
	}
	if !strings.EqualFold(server.Server, "master") {
		t.Errorf("master server value error:%s", server.Server)
	}
	masterValue, err := getRCServerValue(server.Path)
	if err != nil {
		t.Error(err)
	}
	if !strings.EqualFold(masterValue.Domain, zkClient.Domain) {
		t.Error("domain is not correct")
	}
	if !strings.EqualFold(masterValue.IP, zkClient.LocalIP) {
		t.Error("ip is not correct")
	}
	if !strings.EqualFold(masterValue.Online, fmt.Sprintf("%d", server.OnlineTime)) {
		t.Error("online time is not correct")
	}
	if !strings.EqualFold(masterValue.Port, server.Port) {
		t.Error("port is not correct")
	}
	if server.LastPublish <= pt {
		t.Errorf("not publish service provider config:current:%d,last:%d", server.LastPublish, pt)
	}
	if len(server.spServicesMap.data) != spCount {
		t.Errorf("master sp server count is not correct:current:%d,expect:%d",len(server.spServicesMap.data), spCount)
	}
}

func TSlaveServer(server *rcServer, spCount int, t *testing.T) {
	if server.IsMasterServer {
		t.Error("slave bind failed")
	}
	if !strings.EqualFold(server.Server, "slave") {
		t.Errorf("slave server value error:%s", server.Server)
	}
	masterValue, err := getRCServerValue(server.Path)
	if err != nil {
		t.Error(err)
	}
	if !strings.EqualFold(masterValue.Domain, zkClient.Domain) {
		t.Error("domain is not correct")
	}
	if !strings.EqualFold(masterValue.IP, zkClient.LocalIP) {
		t.Error("ip is not correct")
	}
	if !strings.EqualFold(masterValue.Online, fmt.Sprintf("%d", server.OnlineTime)) {
		t.Error("online time is not correct")
	}
	if !strings.EqualFold(masterValue.Port, server.Port) {
		t.Error("port is not correct")
	}
	if len(server.spServicesMap.data) != spCount {
		t.Errorf("slave sp server count is not correct:current:%d,expect:%d", len(server.spServicesMap.data), spCount)
	}
}

func Test_rc(t *testing.T) {

	master := NewRCServer()
	slave := NewRCServer()

	pt := master.LastPublish

	master.Bind()
	slave.Bind()
    
    slave.WatchServiceChange(func(services map[string][]string, err error) {
		slave.BindSPServer(services)
	})

	TMasterServer(master, pt, 0, t)
	TSlaveServer(slave, 0, t)
	spServer := NewSPServer()
	spServer.WatchServiceConfigChange()
    spServer.StartRPC()

	pt = slave.LastPublish
	master.Close()

	time.Sleep(time.Second * 5)
	TMasterServer(slave, pt, 2, t)

}
