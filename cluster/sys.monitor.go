package cluster

import (
	"encoding/json"
	"time"

	"github.com/colinyl/ars/config"
	"github.com/colinyl/lib4go/sysinfo"
	"github.com/colinyl/lib4go/utility"
)

type sysMonitorInfo struct {
	CPU  json.RawMessage `json:"cpu"`
	Mem  json.RawMessage `json:"mem"`
	Disk json.RawMessage `json:"disk"`
}

func GetSysMonitorInfo() (sys *sysMonitorInfo, err error) {
	sys = &sysMonitorInfo{}
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
