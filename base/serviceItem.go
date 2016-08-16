package base

import (
	"errors"
	"sync"
	"sync/atomic"
)

type ServiceGroup struct {
	services []string
	first    int32
	index    int32
	hasUse   int32
	lk       sync.Mutex
}

//GetNext 获取下一个服务
func (g *ServiceGroup) GetNext() (s string, er error) {
	g.lk.Lock()
	defer g.lk.Unlock()
	cindex := g.index % int32(len(g.services))
	cfirst := g.first % int32(len(g.services))

	defer func() {
		g.index++
		if cindex == cfirst {
			atomic.CompareAndSwapInt32(&g.hasUse, 1, 0)
		}
	}()

	if cindex != cfirst || (cindex == cfirst && g.hasUse == 1) {
		s = g.services[cindex]
		return
	}
	er = errors.New("not find services")
	return
}

//ServiceItem 服务信息
type ServiceItem struct {
	services []string
	index    int32
	lk       sync.Mutex
}

//NewServiceItem 创建新的serviceitem
func NewServiceItem(services []string) *ServiceItem {
	return &ServiceItem{services: services}
}

//GetGroup 获取一个可用的服务
func (i *ServiceItem) GetGroup() *ServiceGroup {
	i.lk.Lock()
	defer i.lk.Unlock()
	index := atomic.AddInt32(&i.index, 1)
	cindex := index % int32(len(i.services))

	return &ServiceGroup{services: i.services, first: cindex, index: cindex, hasUse: 1}
}
