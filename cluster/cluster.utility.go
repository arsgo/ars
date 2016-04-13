package cluster

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/colinyl/ars/config"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
	zk "github.com/colinyl/lib4go/zkclient"
)

type zkClientObj struct {
	ZkCli   *zk.ZKCli
	LocalIP string
	Domain  string
	Err     error
	Log     *logger.Logger
}

func (zkClient *zkClientObj) waitZKPathExists(path string, timeout time.Duration, callback func(exists bool)) {
	if zkClient.ZkCli.Exists(path) {
		callback(true)
		return
	}
	callback(false)
	timePiker := time.NewTicker(time.Second * 2)
	timeoutPiker := time.NewTicker(timeout)
CHECKER:
	for {
		select {
		case <-timeoutPiker.C:
			break
		case <-timePiker.C:
			if zkClient.ZkCli.Exists(path) {
				break CHECKER
			}
		}
	}
	callback(zkClient.ZkCli.Exists(path))
}

func (zkClient *zkClientObj) watchZKValueChange(path string, callback func()) {
	changes := make(chan string, 10)
	go zkClient.ZkCli.WatchValue(path, changes)
	go func() {
		for {
			select {
			case <-changes:
				{
					callback()
				}
			}
		}
	}()
}

func (zkClient *zkClientObj) watchZKChildrenPathChange(path string, callback func()) {
	changes := make(chan []string, 10)
	go func() {
		go zkClient.ZkCli.WatchChildren(path, changes)
		for {
			select {
			case <-changes:
				{
					callback()
				}
			}
		}
	}()
}

func (zkClient *zkClientObj) getRCServerValue(path string) (value *RCServerConfig, err error) {
	content, err := zkClient.ZkCli.GetValue(path)
	if err != nil {
		return
	}
	value = &RCServerConfig{}
	err = json.Unmarshal([]byte(content), &value)
	return
}

func (zkClient *zkClientObj) getRCServer(dataMap *utility.DataMap) (servers []*RCServerConfig, err error) {
	path := dataMap.Translate(rcServerRoot)
	rcs, _ := zkClient.ZkCli.GetChildren(path)
	servers = []*RCServerConfig{}
	for _, v := range rcs {
		rcmap := dataMap.Copy()
		rcmap.Set("name", v)
		rcPath := rcmap.Translate(rcServerNodePath)
		config, err := zkClient.getRCServerValue(rcPath)
		if err != nil {
			continue
		}
		servers = append(servers, config)
	}
	return
}

func (zkClient *zkClientObj) getAppConfig(path string) (config *AppConfig, err error) {
	config = &AppConfig{}
	values, err := zkClient.ZkCli.GetValue(path)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(values), &config)
	return
}
func (zkClient *zkClientObj) checkIP(origin string) bool {
	ips := fmt.Sprintf(",%s,", origin)
	llocal := fmt.Sprintf(",%s,", zkClient.LocalIP)
	return strings.Contains(ips, llocal)
}

func (zkClient *zkClientObj) getSPConfig(path string) (svs []*spService, err error) {
	values, err := zkClient.ZkCli.GetValue(path)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(values), &svs)
	return
}

func NewZKClient() *zkClientObj {
	var err error
	client := &zkClientObj{}
	client.Log, err = logger.New("zk client", true)
	client.Domain = config.Get().Domain
	client.LocalIP = utility.GetLocalIP("192.168")
	client.ZkCli, err = zk.New(config.Get().ZKServers, time.Second)
	if err != nil && client.Log != nil {
		client.Log.Error(err)
	}
	return client
}
