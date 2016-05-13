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
	Delete(path string) error
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
	Close()

	//app server..........
	WatchAppTaskChange(callback func(config *AppServerStartupConfig, err error) error)
	GetAppServerStartupConfig(path string) (config *AppServerStartupConfig, err error)
	ResetAppServerSnap(snap string) error
	//rc server...........
	WatchRCServerChange(callback func([]*RCServerItem, error))
	GetRCServerValue(path string) (value *RCServerItem, err error)
	GetAllRCServerValues() (servers []*RCServerItem, err error)
	CreateRCServer(value string) (string, error)
	GetRCServerTasks() (config RCServerTask, err error)
	WatchRCTaskChange(callback func(RCServerTask, error))

	//job server/consumer........
	WatchJobConfigChange(callback func(items map[string]TaskItem, err error))
	GetJobConfig() (items map[string]TaskItem, err error)
	GetJobConsumers(jobName string) (consumers []string)
	CreateJobConsumer(jobName string, value string) (string, error)
	UpdateJobConsumerPath(path string, value string) error
	CloseJobConsumer(path string) error

	//rpc service..........
	WatchRPCServiceChange(callback func(services map[string][]string, err error))
	GetRPCService() (ServiceProviderList, error)

	//sp server........
	WatchServiceProviderChange(changed func()) (err error)
	WatchSPTaskChange(callback func())
	GetAllServiceProviderNamePath() (lst map[string][]string, err error)
	GetServiceTasks() ([]TaskItem, error)
	FilterRPCService(map[string][]string) ([]TaskItem, error)
	PublishRPCServices(map[string]map[string][]string) (err error)
	GetServiceProviderPaths() (lst ServiceProviderList, err error)
	ResetSnap(addr string, snap string) (err error)
	CreateServiceProvider(name string, port string, value string) (string, error)
	CloseServiceProvider(path string) error
}
