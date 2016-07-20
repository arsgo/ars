package cluster

//----------------app server-------------------------------
type JobItem struct {
	Name        string `json:"name"`
	Trigger     string `json:"trigger"`
	Params      string `json:"params"`
	Script      string `json:"script"`
	MinSize     int    `json:"min"`
	MaxSize     int    `json:"max"`
	Concurrency int    `json:"concurrency"`
	Enable      bool   `json:"enable"`
}
type RPCPoolSetting struct {
	MinSize int `json:"min"`
	MaxSize int `json:"max"`
}
type ServerRouteConfig struct {
	Path    string `json:"path"`
	Method  string `json:"method"`
	Script  string `json:"script"`
	Params  string `json:"params"`
	MinSize int    `json:"min"`
	MaxSize int    `json:"max"`
}
type ServerConfig struct {
	Address    string               `json:"address"`
	ServerType string               `json:"type"`
	Routes     []*ServerRouteConfig `json:"routes"`
}

type RootConfig struct {
	Status string         `json:"status"`
	RPC    RPCPoolSetting `json:"rpc"`
	Libs   []string       `json:"libs"`
}
type AppServerStartupConfig struct {
	LocalJobs []JobItem     `json:"jobs"`
	Tasks     []TaskItem    `json:"tasks"`
	Server    *ServerConfig `json:"server"`
	Config    RootConfig    `json:"config"`
}

//---------------------------------------------------------
//----------------app server-------------------------------
type RCServerItem struct {
	Domain  string
	Address string
	Server  string
	Path    string
}

//---------------------------------------------------------
//----------------job server-------------------------------
type JobConsumerValue struct {
	Server string `json:"server"`
}
type ServiceProviderList map[string][]string

//---------------------------------------------------------

//----------------sp server-------------------------------
type TaskItem struct {
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Mode    string `json:"mode"`
	Type    string `json:"type"`
	Method  string `json:"method"`
	Script  string `json:"script"`
	Params  string `json:"params"`
	Trigger string `json:"trigger"`
	MinSize int    `json:"min"`
	MaxSize int    `json:"max"`
}
type ServiceProviderTask struct {
	Config RootConfig `json:"config"`
	Tasks  []TaskItem `json:"tasks"`
}

//---------------------------------------------------------

//----------------rc server-------------------------------
type CrossDoaminAccessItem struct {
	Services []string `json:"services"`
	Type     string   `json:"type"`
	Servers  []string `json:"servers"`
}

type RCServerTask struct {
	CrossDomainAccess map[string]CrossDoaminAccessItem `json:"cross-domain-access"`
	RPCPoolSetting    RPCPoolSetting                   `json:"rpc-pool"`
}

func (c CrossDoaminAccessItem) GetServicesMap() map[string][]string {
	m := make(map[string][]string)
	for _, k := range c.Services {
		m[k] = c.Servers
	}
	return m
}

//---------------------------------------------------------
