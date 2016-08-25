package cluster

import (
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/arsgo/lib4go/concurrent"
	"github.com/arsgo/lib4go/logger"
	"github.com/arsgo/lib4go/utility"
)

const (
	p_varConfig           = "@domain/var/@type/@name"
	p_appTaskRoot         = "@domain/app/config"
	p_appTaskConfig       = "@domain/app/config/@ip"
	p_rcServerTaskConfig  = "@domain/rc/config"
	p_jobTaskConfig       = "@domain/job/config"
	p_spTaskConfig        = "@domain/sp/config"
	p_servicePublishPath  = "@domain/sp/publish"
	p_serviceProviderRoot = "@domain/sp/servers"

	p_appServerPath                    = "@domain/app/api.servers/@ip@port"
	p_serviceProviderPath              = "@domain/sp/servers/@serviceName/@ip@port"
	p_rcServerRoot                     = "@domain/rc/servers"
	p_rcServerClusterClientBase        = "@domain/rc/servers/rc_"
	p_appServerClusterClientPathFormat = "@domain/app/servers/@ip"

	p_MQConsumerClusterClientBase = "@domain/app/mq.consumers/@name/mq_"

	p_jobConsumerClusterClientBase    = "@domain/app/job.consumers/@jobName/job_"
	p_rcServerClusterServerPathFormat = "@domain/rc/servers/@name"

	p_localjobClusterClientBase = "@domain/app/job.locals/@jobName/job_"

	p_jobConsumerNamedRootForamt         = "@domain/app/job.consumers/@jobName"
	p_jobConsumerClusterClientPathFormat = "@domain/app/job.consumers/@jobName/@path"
)

//ClusterClient 集群客户端
type ClusterClient struct {
	domain              string
	domainName          string
	appTaskRoot         string
	appServerTaskPath   string
	rcServerRoot        string
	rcServerConfig      string
	jobConfigPath       string
	spConfigPath        string
	rpcPublishPath      string
	rpcProviderRootPath string
	spServerTaskPath    string
	closeChans          *concurrent.ConcurrentMap
	lastRPCServices     RPCServices
	publishLock         sync.Mutex
	configCache         *concurrent.ConcurrentMap
	handler             IClusterHandler
	Log                 logger.ILogger
	timeout             time.Duration
	dataMap             utility.DataMap
	IP                  string
}

func NewClusterClient(domain string, ip string, handler IClusterHandler, loggerName string) (client *ClusterClient, err error) {
	client = &ClusterClient{configCache: concurrent.NewConcurrentMap()}
	client.domain = "/" + strings.Trim(strings.Replace(domain, ".", "/", -1), "/")
	client.domainName = "@" + strings.Replace(strings.Trim(client.domain, "/"), "/", ".", -1)
	client.IP = ip
	client.dataMap = utility.NewDataMap()
	client.dataMap.Set("domain", client.domain)
	client.dataMap.Set("ip", client.IP)
	client.closeChans = concurrent.NewConcurrentMap()
	client.appServerTaskPath = client.dataMap.Translate(p_appTaskConfig)
	client.rcServerRoot = client.dataMap.Translate(p_rcServerRoot)
	client.appTaskRoot = client.dataMap.Translate(p_appTaskRoot)
	client.rcServerConfig = client.dataMap.Translate(p_rcServerTaskConfig)
	client.spConfigPath = client.dataMap.Translate(p_spTaskConfig)
	client.rpcPublishPath = client.dataMap.Translate(p_servicePublishPath)
	client.rpcProviderRootPath = client.dataMap.Translate(p_serviceProviderRoot)
	client.jobConfigPath = client.dataMap.Translate(p_jobTaskConfig)
	client.spServerTaskPath = client.dataMap.Translate(p_spTaskConfig)
	client.Log, err = logger.Get(loggerName)
	client.timeout = time.Hour * 10000
	client.handler = handler
	client.handler.Open()
	return
}
func (client *ClusterClient) makeCloseChan() chan int {
	closeChan := make(chan int, 1)
	client.closeChans.Set(utility.GetGUID(), closeChan)
	return closeChan
}
func (client *ClusterClient) GetDomainName() string {
	return client.domainName
}
func (client *ClusterClient) GetHandler() IClusterHandler {
	return client.handler
}

//WaitForConnected 监控是否已链接到当前服务器
func (client *ClusterClient) WaitForConnected() bool {
	return client.handler.WaitForConnected()
}

//WaitForDisconnected 监控是否已链接到当前服务器
func (client *ClusterClient) WaitForDisconnected() bool {
	return client.handler.WaitForDisconnected()
}

//Reconnect 重新连接到服务器
func (client *ClusterClient) Reconnect() error {
	return client.handler.Reconnect()
}

//Close 关闭当前集群客户端
func (client *ClusterClient) Close() {
	all := client.closeChans.GetAll()
	for _, v := range all {
		ch := v.(chan int)
		ch <- 1
	}
	client.handler.Close()
}

//recover 从异常中恢复
func (client *ClusterClient) recover() {
	if r := recover(); r != nil {
		client.Log.Fatal(r, string(debug.Stack()))
	}
}
