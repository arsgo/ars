package cluster

import (
	"encoding/json"
	"log"
	"time"

	"github.com/colinyl/ars/config"
	"github.com/colinyl/lib4go/utility"
	zk "github.com/colinyl/lib4go/zkclient"
)

type zkClientObj struct {
	ZkCli   *zk.ZKCli
	LocalIP string
	Domain  string
	Err     error
}

func waitZKPathExists(path string, timeout time.Duration, callback func(exists bool)) {
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

func watchZKValueChange(path string, callback func()) {
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

func watchZKChildrenPathChange(path string, callback func()) {
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

func getRCServerValue(path string) (value *RCServerConfig, err error) {
	content, err := zkClient.ZkCli.GetValue(path)
	if err != nil {
		return
	}
	value = &RCServerConfig{}
	err = json.Unmarshal([]byte(content), &value)
	return
}

func getRCServer(dataMap *utility.DataMap) (servers []*RCServerConfig, err error) {
	path := dataMap.Translate(rcServerRoot)
	rcs, _ := zkClient.ZkCli.GetChildren(path)
	servers = []*RCServerConfig{}
	for _, v := range rcs {
		rcmap := dataMap.Copy()
		rcmap.Set("name", v)
		rcPath := rcmap.Translate(rcServerNodePath)
		config, err := getRCServerValue(rcPath)
		if err != nil {
			continue
		}
		servers = append(servers, config)
	}
	return
}

func getAppConfig(path string) (config *AppConfig, err error) {
	config = &AppConfig{}
	values, err := zkClient.ZkCli.GetValue(path)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(values), &config)
	return
}

var zkClient *zkClientObj

func init() {
	var err error
	zkClient = &zkClientObj{}
	zkClient.Domain = config.Get().Domain
	zkClient.LocalIP = utility.GetLocalIP("192.168")
	zkClient.ZkCli, err = zk.New(config.Get().ZKServers, time.Second)
	if err != nil {
		log.Println(err)
	}
}
