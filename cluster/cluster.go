package cluster

import "time"

//IClusterHandler 集群管理处理程序，用于处理与集群管理器如(zk,etcd)等之前的通信
type IClusterHandler interface {
	Open()
	Exists(path ...string) (string, bool)
	CreateNode(path string, data string) error
	CreateSeqNode(path string, data string) (string, error)
	CreateTmpNode(path string, data string) (string, error)
	GetValue(path string) (string, error)
	UpdateValue(path string, value string) error
	GetChildren(path string) ([]string, error)
	WatchValue(path string, data chan string) error
	WatchChildren(path string, data chan []string) error
	RemoveWatchValue(path string)
	RemoveWatchChildren(path string)
	WaitForConnected() bool
	WaitForDisconnected() bool
	Delete(path string) error
	Reconnect() error
	Close()
}

//IClusterClient 集群客户端处理程序，提供appserver,rcserver,spserver与集群之前的交互操作
type IClusterClient interface {
	//base.............
	GetHandler() IClusterHandler
	GetDomainName() string
	WaitClusterPathExists(path string, timeout time.Duration, callback func(path string, exists bool))
	WatchClusterValueChange(path string, callback func())
	WatchClusterChildrenChange(path string, callback func())
	GetSourceConfig(typeName string, name string) (config string, err error)
	GetMQConfig(name string) (string, error)
	GetElasticConfig(name string) (string, error)
	GetDBConfig(name string) (string, error)
	GetServiceFullPath(name string) string
	CreateNode(path string, value string) (string, error)
	UpdateNode(path string, value string) error
	CloseNode(path string) error
	SetNode(path string, value string) error
	WaitForConnected() bool
	WaitForDisconnected() bool
	Reconnect() error
	Close()

	//app server..........
	WatchAppTaskChange(callback func(config *AppServerTask, err error) error)
	GetCurrentAppServerTask() (config *AppServerTask, err error)
	GetAppServerTaskNames() ([]string, error)
	GetAppServerTask(name string) (config *AppServerTask, err error)
	UpdateAppServerTask(name string, config *AppServerTask) (err error)
	CreateAppServer(port string, snap string) (string, error)
	CloseAppServer(path string) error

	//rc server...........
	WatchRCServerChange(callback func([]*RCServerItem, error))
	GetRCServerValue(path string) (value *RCServerItem, err error)
	GetAllRCServers() (servers []*RCServerItem, err error)
	CreateRCServer(value string) (string, error)
	CloseRCServer(path string) error
	GetRCServerTask() (config RCServerTask, err error)
	UpdateRCServerTask(config RCServerTask) error
	WatchRCTaskChange(callback func(RCServerTask, error))

	//job server........
	WatchJobConfigChange(callback func(items map[string]JobItem, err error))
	GetJobTask() (items map[string]JobItem, err error)
	UpdateJobTask(jobName string, items map[string]JobItem) (err error)

	//job consumer
	GetJobConsumers(jobName string) (consumers []string)
	CreateJobConsumer(jobName string, value string) (string, error)
	CreateLocalJob(jobName string, value string) (string, error)
	//mq consumer
	CreateMQConsumer(name string, value string) (string, error)

	//rpc service..........
	WatchRPCServiceChange(callback func(services map[string][]string, err error))
	GetPublishServices() (RPCServices, error)
	GetSPServerServices() (lst RPCServices, err error)
	PublishServices(services RPCServices) (err error)

	//sp server........
	WatchSPServerChange(changed func(RPCServices, error)) (err error)
	WatchSPTaskChange(callback func(SPServerTask, error))
	GetAllSPServers() (lst map[string][]string, err error)
	GetSPServerTask(ip string) (SPServerTask, error)
	UpdateSPServerTask(task SPServerTask) (err error)
	GetLocalServices(map[string][]string) ([]TaskItem, error)
	CreateSPServer(name string, port string, value string) (string, error)
	CloseSPServer(path string) error
}
