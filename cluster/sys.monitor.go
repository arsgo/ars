package cluster

import (
	"errors"
	"strings"

	"github.com/colinyl/ars/scheduler"
	"github.com/colinyl/ars/sys"
)

type logHandler interface {
	Infof(format string, a ...interface{})
	Error(content interface{})
	Info(content string)
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
	Cpu  *monitorItemConfig `json:"cpu"`
	Mem  *monitorItemConfig `json:"mem"`
	Disk *monitorItemConfig `json:"disk"`
}

type serverMonitor struct {
	sch    *scheduler.Scheduler
	hander SourceHandler
	Log    logHandler
}

func NewMonitor(h SourceHandler, log logHandler) *serverMonitor {
	return &serverMonitor{sch: scheduler.NewScheduler(), hander: h, Log: log}
}

func (s *serverMonitor) Bind(c *monitorConfig) (err error) {
	s.sch.Stop()
	if c == nil {
		return nil
	}
	if c.Cpu != nil && !strings.EqualFold(c.Cpu.Trigger, "") {
		content, err := s.checkParams(c.Cpu)
		if err == nil {
			c.Cpu.content = content
			s.sch.AddTask(c.Cpu.Trigger, scheduler.NewTask(c.Cpu, func(obj interface{}) {
				cpu := obj.(*monitorItemConfig)
				handler, err := getMonitorHandler(cpu.Source.TypeName, cpu.content)
				if err == nil {
					s.Log.Info(">send cpu info")
					err = handler.Send(cpu.Source.Param, sys.GetCPU())
					if err != nil {
						s.Log.Error(err)
					}
				} else {
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
				handler, err := getMonitorHandler(mem.Source.TypeName, mem.content)
				if err == nil {
					s.Log.Info(">send mem info")
					err = handler.Send(mem.Source.Param, sys.GetMemory())
					if err != nil {
						s.Log.Error(err)
					}
					
				} else {
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
				handler, err := getMonitorHandler(disk.Source.TypeName, disk.content)
				if err == nil {
					s.Log.Info(">send disk info")
					err = handler.Send(disk.Source.Param, sys.GetDisk())
					if err != nil {
						s.Log.Error(err)
					}
				} else {
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