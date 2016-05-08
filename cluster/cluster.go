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

	//app server..........
	WatchAppTaskChange(callback func(config *AppServerStartupConfig, err error) error)
	GetAppServerStartupConfig(path string) (config *AppServerStartupConfig, err error)
	ResetAppServerSnap(snap string) error
	//rc server...........
	WatchRCServerChange(callback func([]*RCServerItem, error))
	GetRCServerValue(path string) (value *RCServerItem, err error)
	GetAllRCServerValues() (servers []*RCServerItem, err error)
	CreateRCServer(value string) (string, error)

	//job server/consumer........
	WatchJobConfigChange(callback func(items *JobItems, err error))
	GetJobConfig(path string) (items *JobItems, err error)
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
	GeServiceTasks() ([]TaskItem, error)
	PublishRPCServices() (err error)
	GetServiceProviderPaths() (lst ServiceProviderList, err error)
	ResetSnap(addr string, snap string) (err error)
	CreateServiceProvider(name string, port string, value string) (string, error)
	CloseServiceProvider(path string) error
}
