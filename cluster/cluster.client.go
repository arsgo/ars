package cluster

import (
	"strings"
	"time"

	"github.com/colinyl/lib4go/concurrent"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

const (
	p_varConfig           = "@domain/var/@type/@name"
	p_appTaskConfig       = "@domain/app/config/@ip"
	p_rcServerTaskConfig  = "@domain/rc/config"
	p_jobTaskConfig       = "@domain/job/config"
	p_spTaskConfig        = "@domain/sp/config"
	p_servicePublishPath  = "@domain/sp/publish"
	p_serviceProviderRoot = "@domain/sp/servers"

	p_appServerPath                    = "@domain/app/servers/@ip"
	p_serviceProviderPath              = "@domain/sp/servers/@serviceName/@ip@port"
	p_rcServerRoot                     = "@domain/rc/servers"
	p_rcServerClusterClientBase        = "@domain/rc/servers/rc_"
	p_appServerClusterClientPathFormat = "@domain/app/servers/@ip"

	p_jobConsumerClusterClientBase    = "@domain/job/servers/@jobName/job_"
	p_rcServerClusterServerPathFormat = "@domain/rc/servers/@name"

	p_jobConsumerNamedRootForamt         = "@domain/job/servers/@jobName"
	p_jobConsumerClusterClientPathFormat = "@domain/job/servers/@jobName/@path"
)

//ClusterClient 集群客户端
type ClusterClient struct {
	domain              string
	domainPath          string
	appServerTaskPath   string
	rcServerRoot        string
	rcServerConfig      string
	jobConfigPath       string
	spConfigPath        string
	rpcPublishPath      string
	rpcProviderRootPath string
	appServerPath       string
	spServerTaskPath    string
	configCache         concurrent.ConcurrentMap
	handler             IClusterHandler
	Log                 logger.ILogger
	timeout             time.Duration
	dataMap             utility.DataMap
	IP                  string
}

func NewClusterClient(domain string, ip string, handler IClusterHandler, loggerName string) (client *ClusterClient, err error) {
	client = &ClusterClient{configCache: concurrent.NewConcurrentMap()}
	client.domain = "/" + strings.TrimLeft(strings.Replace(domain, ".", "/", -1), "/")
	client.domainPath = "@" + strings.Replace(strings.TrimLeft(client.domain, "/"), "/", ".", -1)
	client.IP = ip
	client.dataMap = utility.NewDataMap()
	client.dataMap.Set("domain", client.domain)
	client.dataMap.Set("ip", client.IP)
	client.appServerTaskPath = client.dataMap.Translate(p_appTaskConfig)
	client.rcServerRoot = client.dataMap.Translate(p_rcServerRoot)
	client.rcServerConfig = client.dataMap.Translate(p_rcServerTaskConfig)
	client.spConfigPath = client.dataMap.Translate(p_spTaskConfig)
	client.rpcPublishPath = client.dataMap.Translate(p_servicePublishPath)
	client.rpcProviderRootPath = client.dataMap.Translate(p_serviceProviderRoot)
	client.jobConfigPath = client.dataMap.Translate(p_jobTaskConfig)
	client.appServerPath = client.dataMap.Translate(p_appServerPath)
	client.spServerTaskPath = client.dataMap.Translate(p_spTaskConfig)
	client.Log, err = logger.Get(loggerName, true)
	client.timeout = time.Hour * 10000
	client.handler = handler
	return
}

//WaitClusterPathExists  等待集群中的指定配置出现,不存在时持续等待
func (client *ClusterClient) WaitClusterPathExists(path string, timeout time.Duration, callback func(exists bool)) {
	if client.handler.Exists(path) {
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
			if client.handler.Exists(path) {
				break CHECKER
			}
		}
	}
	callback(client.handler.Exists(path))
}

//WatchClusterValueChange 等待集群指定路径的值的变化
func (client *ClusterClient) WatchClusterValueChange(path string, callback func()) {
	changes := make(chan string, 10)
	go func() {
		defer client.recover()
		client.handler.WatchValue(path, changes)
	}()
	go func() {
		for {
			select {
			case <-changes:
				{
					defer client.recover()
					callback()
				}
			}
		}
	}()
}

//WatchClusterChildrenChange 监控集群指定路径的子节点变化
func (client *ClusterClient) WatchClusterChildrenChange(path string, callback func()) {
	changes := make(chan []string, 10)
	go func() {

		go func() {
			defer client.recover()
			client.handler.WatchChildren(path, changes)
		}()
		for {
			select {
			case <-changes:
				{
					defer client.recover()
					callback()
				}
			}
		}
	}()
}

//WatchConnected 监控是否已链接到当前服务器
func (client *ClusterClient) WatchConnected() bool {
	return client.handler.WatchConnected()
}

//Close 关闭当前集群客户端
func (client *ClusterClient) Close() {
	client.handler.Close()
}
