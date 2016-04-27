package sys

import (
	"encoding/json"
	"runtime"
	"time"

	"github.com/colinyl/lib4go/utility"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type stat struct {
	Data interface{} `json:"data"`
	IP   string      `json:"ip"`
	Time int64       `json:"time"`
}

func GetMemory() string {
	v, _ := mem.VirtualMemory()
	data, _ := json.Marshal(&stat{v, utility.GetLocalIP("192"), time.Now().Unix()})
	return string(data)
}
func GetCPU() string {
	v, _ := cpu.Times(true)
	data, _ := json.Marshal(&stat{v, utility.GetLocalIP("192"), time.Now().Unix()})
	return string(data)
}

func GetDisk() string {
	var stats []*disk.UsageStat
	if runtime.GOOS == "windows" {
		v, _ := disk.Partitions(true)
		for _, p := range v {
			s, _ := disk.Usage(p.Device)
			stats = append(stats, s)
		}
	} else {
		s, _ := disk.Usage("/")
		stats = append(stats, s)
	}
	data, _ := json.Marshal(&stat{stats, utility.GetLocalIP("192"), time.Now().Unix()})
	return string(data)
}
