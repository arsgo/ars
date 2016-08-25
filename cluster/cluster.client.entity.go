package cluster

import "strings"

const (
	SERVER_MASTER  = "master"
	SERVER_SLAVE   = "slave"
	SERVER_UNKNOWN = "unknown"
)

//----------------app server-------------------------------
type JobItem struct {
	Name        string `json:"name"`
	Trigger     string `json:"trigger"`
	Params      string `json:"params"`
	Script      string `json:"script"`
	MinSize     int    `json:"min"`
	MaxSize     int    `json:"max"`
	Concurrency int    `json:"concurrency"`
	Disable     bool   `json:"disable"`
}
type RPCPoolSetting struct {
	MinSize int `json:"min"`
	MaxSize int `json:"max"`
}
type ServerRouteConfig struct {
	Path     string `json:"path"`
	Method   string `json:"method"`
	Script   string `json:"script"`
	Params   string `json:"params"`
	Encoding string `json:"chaset"`
	MinSize  int    `json:"min"`
	MaxSize  int    `json:"max"`
}
type ServerConfig struct {
	Address    string               `json:"address"`
	ServerType string               `json:"type"`
	Disable    bool                 `json:"disable"`
	Routes     []*ServerRouteConfig `json:"routes"`
}

type RootConfig struct {
	Status      string         `json:"status"`      //服务器状态，用于启用或停用服务器所有服务(暂未提供)
	RPC         RPCPoolSetting `json:"rpc"`         //rpc 缓存池配置
	Libs        []string       `json:"libs"`        //脚本库，脚本引擎加载时用于设置环境路径
	DisableRPC  bool           `json:"disable-rpc"` //禁用rpc后，不再缓存与监控RPC Server的变化
	SnapRefresh int            `json:"refresh"`     //快照刷新时间，单位秒,不能低于60秒
}
type AppServerTask struct {
	Tasks  []TaskItem    `json:"tasks"`
	Server *ServerConfig `json:"api"`
	Config RootConfig    `json:"config"`
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
	Disable bool   `json:"disable"`
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
	Disable  bool     `json:"disable"`
}

type RCServerTask struct {
	SnapRefresh       int                              `json:"refresh"` //快照刷新时间，单位秒,不能低于60秒
	CrossDomainAccess map[string]CrossDoaminAccessItem `json:"cross-domain"`
	RPCPoolSetting    RPCPoolSetting                   `json:"rpc"`
	Jobs              []JobItem                        `json:"jobs"`
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
