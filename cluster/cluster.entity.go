package cluster

//----------------app server-------------------------------
type TaskConfig struct {
	Trigger string `json:"trigger"`
	Script  string `json:"script"`
}
type ServerRouteConfig struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Script string `json:"script"`
}
type ServerConfig struct {
	ServerType string               `json:"type"`
	Routes     []*ServerRouteConfig `json:"routes"`
}
type AppServerStartupConfig struct {
	Status   string        `json:"status"`
	Tasks    []*TaskConfig `json:"tasks"`
	JobNames []string      `json:"jobs"`
	Jobs     []TaskItem
	Server   *ServerConfig `json:"server"`
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
	Address string `json:"address"`
}
type ServiceProviderList map[string][]string

//---------------------------------------------------------

//----------------sp server-------------------------------
type TaskItem struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Mode        string `json:"mode"`
	Type        string `json:"type"`
	Method      string `json:"method"`
	Script      string `json:"script"`
	Params      string `json:"params"`
	Trigger     string `json:"trigger"`
	Concurrency int    `json:"concurrency"`
}

//---------------------------------------------------------

//----------------rc server-------------------------------
type CrossDoaminAccessItem struct {
	Services []string `json:"service"`
	Type     string   `json:"type"`
	Servers  []string `json:"servers"`
}
type RCServerTask struct {
	CrossDomainAccess map[string]CrossDoaminAccessItem `json:"cross-domain-access"`
}

//---------------------------------------------------------
