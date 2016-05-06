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

const (
	varConfigPath = "@domain/var/@type/@name"
)

type clusterClient struct {
	ZkCli   *zk.ZKCli
	IP      string
	Domain  string
	Err     error
	Log     *logger.Logger
	dataMap utility.DataMap
}

func (client *clusterClient) waitZKPathExists(path string, timeout time.Duration, callback func(exists bool)) {
	if client.ZkCli.Exists(path) {
		callback(true)
		return
	}
	callback(false)
	timePiker := time.NewTicker(time.Second * 2)
	timeoutPiker := time.NewTicker(timeout)
	defer func() {
		timeoutPiker.Stop()
	}()
CHECKER:
	for {
		select {
		case <-timeoutPiker.C:
			break
		case <-timePiker.C:
			if client.ZkCli.Exists(path) {
				break CHECKER
			}
		}
	}
	callback(client.ZkCli.Exists(path))
}

func (client *clusterClient) watchZKValueChange(path string, callback func()) {
	changes := make(chan string, 10)
	go client.ZkCli.WatchValue(path, changes)
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

func (client *clusterClient) watchZKChildrenPathChange(path string, callback func()) {
	changes := make(chan []string, 10)
	go func() {
		go client.ZkCli.WatchChildren(path, changes)
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

func (client *clusterClient) getRCServerValue(path string) (value *RCServerConfig, err error) {
	content, err := client.ZkCli.GetValue(path)
	if err != nil {
		return
	}
	value = &RCServerConfig{}
	err = json.Unmarshal([]byte(content), &value)
	return
}

func (client *clusterClient) getRCServer(dataMap utility.DataMap) (servers []*RCServerConfig, err error) {
	path := dataMap.Translate(rcServerRoot)
	rcs, _ := client.ZkCli.GetChildren(path)
	servers = []*RCServerConfig{}
	for _, v := range rcs {
		rcmap := dataMap.Copy()
		rcmap.Set("name", v)
		rcPath := rcmap.Translate(rcServerNodePath)
		config, err := client.getRCServerValue(rcPath)
		if err != nil {
			continue
		}
		servers = append(servers, config)
	}
	return
}

func (client *clusterClient) getAppConfig(path string) (config *AppConfig, err error) {
	config = &AppConfig{}
	values, err := client.ZkCli.GetValue(path)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(values), &config)
	return
}
func (client *clusterClient) checkIP(origin string) bool {
	if strings.EqualFold(origin, "*") {
		return true
	}
	ips := fmt.Sprintf(",%s,", origin)
	llocal := fmt.Sprintf(",%s,", client.IP)
	return strings.Contains(ips, llocal)
}
func (client *clusterClient) GetSourceConfig(typeName string, name string) (config string, err error) {
	dataMap := client.dataMap.Copy()
	dataMap.Set("type", typeName)
	dataMap.Set("name", name)
	values, err := client.ZkCli.GetValue(dataMap.Translate(varConfigPath))
	if err != nil {
		fmt.Println(dataMap.Translate(varConfigPath))
		return
	}
	config = string(values)
	return
}
func (client *clusterClient) getSPConfig(path string) (svs []spService, err error) {
	values, err := client.ZkCli.GetValue(path)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(values), &svs)
	return
}
func (client *clusterClient) GetMQConfig(name string) (string, error) {
	return client.GetSourceConfig("mq", name)
}
func (client *clusterClient) GetElasticConfig(name string) (string, error) {
	return client.GetSourceConfig("elastic", name)
}

func NewClusterClient() *clusterClient {
	var err error
	client := &clusterClient{}
	client.Log, err = logger.New("zk client", true)
	client.Domain = config.Get().Domain
	client.IP = config.Get().IP
	client.ZkCli, err = zk.New(config.Get().ZKServers, time.Second)
	client.dataMap = utility.NewDataMap()
	client.dataMap.Set("ip", client.IP)
	client.dataMap.Set("domain", client.Domain)
	if err != nil && client.Log != nil {
		client.Log.Error(err)
	}
	return client
}
