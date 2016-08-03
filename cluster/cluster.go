package cluster

import "time"

type IClusterHandler interface {
	Exists(path string) bool
	CreatePath(path string, data string) error
	CreateSeqNode(path string, data string) (string, error)
	CreateTmpNode(path string, data string) (string, error)
	GetValue(path string) (string, error)
	GetChildren(path string) ([]string, error)
	UpdateValue(path string, value string) error
	WatchValue(path string, data chan string) error
	WatchChildren(path string, data chan []string) error
	WaitForConnected() bool
	WaitForDisconnected() bool
	Delete(path string) error
	Reconnect() error
	Close()
}
type IClusterClient interface {
	//base.............
	WaitClusterPathExists(path string, timeout time.Duration, callback func(exists bool))
	WatchClusterValueChange(path string, callback func())
	WatchClusterChildrenChange(path string, callback func())
	GetSourceConfig(typeName string, name string) (config string, err error)
	GetMQConfig(name string) (string, error)
	GetElasticConfig(name string) (string, error)
	GetDBConfig(name string) (string, error)
	GetServiceFullPath(name string) string
	WaitForConnected() bool
	WaitForDisconnected() bool
	Reconnect() error
	Close()

	//app server..........
	WatchAppTaskChange(callback func(config *AppServerStartupConfig, err error) error)
	GetCurrentAppServerTask() (config *AppServerStartupConfig, err error)
	GetAppServerTask(name string) (config *AppServerStartupConfig, err error)
	UpdateAppServerTask(name string, config *AppServerStartupConfig) (err error)
	UpdateAppServerSnap(snap string) error
	CloseAppServer() error

	//rc server...........
	WatchRCServerChange(callback func([]*RCServerItem, error))
	GetRCServerValue(path string) (value *RCServerItem, err error)
	GetAllRCServerValues() (servers []*RCServerItem, err error)
	CreateRCServer(value string) (string, error)
	CloseRCServer(path string) error
	GetRCServerTask() (config RCServerTask, err error)
	UpdateRCServerTask(config RCServerTask) error
	WatchRCTaskChange(callback func(RCServerTask, error))

	//job server/consumer........
	WatchJobConfigChange(callback func(items map[string]JobItem, err error))
	GetJobTask() (items map[string]JobItem, err error)
	UpdateJobTask(jobName string, items map[string]JobItem) (err error)
	GetJobConsumers(jobName string) (consumers []string)
	CreateJobConsumer(jobName string, value string) (string, error)
	UpdateJobConsumerPath(path string, value string) error
	CloseJobConsumer(path string) error

	//rpc service..........
	WatchRPCServiceChange(callback func(services map[string][]string, err error))
	GetRPCService() (ServiceProviderList, error)

	//sp server........
	WatchServiceProviderChange(changed func(ServiceProviderList, error)) (err error)
	WatchSPTaskChange(callback func())
	GetAllServiceProviderNamePath() (lst map[string][]string, err error)
	GetSPServerTask() (ServiceProviderTask, error)
	UpdateSPServerTask(task ServiceProviderTask) (err error)
	FilterRPCService(map[string][]string) ([]TaskItem, error)
	PublishRPCServices(services ServiceProviderList) (err error)
	GetServiceProviders() (lst ServiceProviderList, err error)
	UpdateSnap(addr string, snap string) (err error)
	CreateServiceProvider(name string, port string, value string) (string, error)
	CloseServiceProvider(path string) error
}
