package cluster

import (
	"errors"
	"fmt"
	"strings"

	"github.com/colinyl/lib4go/scheduler"
	"github.com/colinyl/lib4go/sysinfo"
)

type logHandler interface {
	Infof(format string, a ...interface{})
	Error(content interface{})
	Info(content string)
}
type requestHandler interface {
	request(string, string) []string
}

type SourceHandler interface {
	GetSourceConfig(string, string) (string, error)
}
type monitorSourceConfig struct {
	TypeName string `json:"type"`
	Name     string `json:"name"`
	Param    string `json:"param"`
}
type monitorItemConfig struct {
	Source  *monitorSourceConfig `json:"source"`
	Trigger string               `json:"trigger"`
	content string
}

type monitorConfig struct {
	CPU  *monitorItemConfig `json:"cpu"`
	Mem  *monitorItemConfig `json:"mem"`
	Disk *monitorItemConfig `json:"disk"`
}

type serverMonitor struct {
	sch    *scheduler.Scheduler
	hander SourceHandler
	Log    logHandler
	r      requestHandler
}

func NewMonitor(h SourceHandler, log logHandler, r requestHandler) *serverMonitor {
	return &serverMonitor{sch: scheduler.NewScheduler(), hander: h, Log: log, r: r}
}

func (s *serverMonitor) Bind(c *monitorConfig) (err error) {
	s.sch.Stop()
	if c == nil {
		return nil
	}
	if c.CPU != nil && !strings.EqualFold(c.CPU.Trigger, "") {
		content, err := s.checkParams(c.CPU)
		if err == nil {
			c.CPU.content = content
			s.sch.AddTask(c.CPU.Trigger, scheduler.NewTask(c.CPU, func(obj interface{}) {
				cpu := obj.(*monitorItemConfig)
				s.Log.Info("->send cpu")
				//fmt.Println("->send cpu")
				err := StaticSendMonitor(cpu.Source.TypeName, cpu.content, cpu.Source.Param, sysinfo.GetCPU())
				//err := mqservice.StaticSend(cpu.Source.Param, sys.GetCPU())
				fmt.Println(err)
				if err != nil {
					s.Log.Error(err)
				}
			}))
		} else {
			s.Log.Error(err)
		}
	}
	if c.Mem != nil && !strings.EqualFold(c.Mem.Trigger, "") {
		content, err := s.checkParams(c.Mem)
		if err == nil {
			c.Mem.content = content
			s.sch.AddTask(c.Mem.Trigger, scheduler.NewTask(c.Mem, func(obj interface{}) {
				mem := obj.(*monitorItemConfig)
				//	fmt.Println("->send mem")
				s.Log.Info("->send mem")
				err := StaticSendMonitor(mem.Source.TypeName, mem.content, mem.Source.Param, sysinfo.GetMemory())
				//err := mqservice.StaticSend(mem.Source.Param, sys.GetMemory())
				if err != nil {
					s.Log.Error(err)
				}
			}))
		} else {
			s.Log.Error(err)
		}
	}
	if c.Disk != nil && !strings.EqualFold(c.Disk.Trigger, "") {
		content, err := s.checkParams(c.Disk)
		if err == nil {
			c.Disk.content = content
			s.sch.AddTask(c.Disk.Trigger, scheduler.NewTask(c.Disk, func(obj interface{}) {
				disk := obj.(*monitorItemConfig)
				//fmt.Println("->send disk")
				s.Log.Info("->send disk")
				err := StaticSendMonitor(disk.Source.TypeName, disk.content, disk.Source.Param, sysinfo.GetDisk())
				//err := mqservice.StaticSend(disk.Source.Param, sys.GetDisk())
				if err != nil {
					s.Log.Error(err)
				}
			}))
		} else {
			s.Log.Error(err)
		}
	}
	s.sch.Start()
	return
}
func (s *serverMonitor) checkParams(c *monitorItemConfig) (content string, err error) {
	err = errors.New("input args error")
	if c.Source == nil {
		return
	}
	if strings.EqualFold(c.Source.TypeName, "") ||
		strings.EqualFold(c.Source.Name, "") {
		return
	}
	content, err = s.hander.GetSourceConfig(c.Source.TypeName, c.Source.Name)
	if err != nil {
		return
	}
	return
}
