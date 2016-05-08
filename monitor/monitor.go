package monitor

import (
	"encoding/json"
	"time"

	"github.com/colinyl/ars.bak/config"
	"github.com/colinyl/lib4go/sysinfo"
	"github.com/colinyl/lib4go/utility"
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
	baseMap := make(map[string]interface{})
	baseMap["ip"] = config.Get().IP
	baseMap["timestamp"] = time.Now().Format("20060102150405")
	sys.CPU, err = json.Marshal(utility.MergeMaps(baseMap, sysinfo.GetCPU()))
	if err != nil {
		return
	}
	sys.Mem, err = json.Marshal(utility.MergeMaps(baseMap, sysinfo.GetMemory()))
	if err != nil {
		return
	}
	sys.Disk, err = json.Marshal(utility.MergeMaps(baseMap, sysinfo.GetDisk()))
	if err != nil {
		return
	}
	return
}
