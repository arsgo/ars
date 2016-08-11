package cluster

import "strings"

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
	Status      string         `json:"status"`
	RPC         RPCPoolSetting `json:"rpc"`
	Libs        []string       `json:"libs"`
	DisableRPC  bool           `json:"disable-rpc"`
	SnapRefresh int            `json:"refresh"`
}
type AppServerTask struct {
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
type RPCServices map[string][]string

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
type SPServerTask struct {
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
	CrossDomainAccess map[string]CrossDoaminAccessItem `json:"cross-domain"`
	RPCPoolSetting    RPCPoolSetting                   `json:"rpc"`
}

func (c CrossDoaminAccessItem) GetServicesMap(domain string) RPCServices {
	m := make(RPCServices)
	for _, k := range c.Services {
		m[k+"@"+domain] = []string{}
	}
	return m
}
func (c RPCServices) Clone() RPCServices {
	mp := make(RPCServices)
	for k, v := range c {
		vs := make([]string, 0, len(v))
		for _, p := range v {
			vs = append(vs, p)
		}
		mp[k] = vs
	}
	return mp
}

func (c RPCServices) Equal(input RPCServices) bool {
	if input == nil {
		return false
	}
	if len(c) != len(input) {
		return false
	}
	for i, v := range c {
		if _, ok := input[i]; !ok {
			return false
		}
		if len(v) != len(input[i]) {
			return false
		}
		for _, m := range v {
			hasExist := true
			for _, n := range input[i] {
				if !strings.EqualFold(m, n) {
					hasExist = false
					break
				}
			}
			if !hasExist {
				return false
			}
		}
	}
	return true
}

//---------------------------------------------------------
