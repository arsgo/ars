package base

import (
	"encoding/json"

	"github.com/colinyl/lib4go/sysinfo"
)

//SysMonitorInfo 系统信息
type SysMonitorInfo struct {
	CPU  json.RawMessage `json:"cpu"`
	Mem  json.RawMessage `json:"mem"`
	Disk json.RawMessage `json:"disk"`
}

//GetSysMonitorInfo 获取系统信息
func GetSysMonitorInfo() (sys *SysMonitorInfo, err error) {
	sys = &SysMonitorInfo{}
	sys.CPU, err = json.Marshal(sysinfo.GetCPU())
	if err != nil {
		return
	}
	sys.Mem, err = json.Marshal(sysinfo.GetMemory())
	if err != nil {
		return
	}
	sys.Disk, err = json.Marshal(sysinfo.GetDisk())
	if err != nil {
		return
	}
	return
}
