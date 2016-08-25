package base

import (
	"fmt"
	"sync/atomic"

	"github.com/arsgo/lib4go/concurrent"
)

type ICollector interface {
	Success(...interface{})
	Failed(...interface{})
	Error(...interface{})
	Juge(v bool, name ...interface{})
}
type Execution struct {
	Success int32 `json:"success"`
	Failed  int32 `json:"failed"`
	Error   int32 `json:"error"`
}

//Collector 采集器
type Collector struct {
	data *concurrent.ConcurrentMap
}

func NewCollector() *Collector {
	r := &Collector{}
	r.data = concurrent.NewConcurrentMap()
	return r
}
func (r *Collector) getExecution(name ...interface{}) (d *Execution, err error) {
	exec, err := r.data.GetOrAdd(fmt.Sprint(name...), func(p ...interface{}) (interface{}, error) {
		return &Execution{}, nil
	})
	if err != nil {
		return
	}
	d = exec.(*Execution)
	return
}
func (r *Collector) Success(name ...interface{}) {
	if data, err := r.getExecution(name...); err == nil {
		atomic.AddInt32(&data.Success, 1)
	}

}
func (r *Collector) Failed(name ...interface{}) {
	if data, err := r.getExecution(name...); err == nil {
		atomic.AddInt32(&data.Failed, 1)
	}
}
func (r *Collector) Error(name ...interface{}) {
	if data, err := r.getExecution(name...); err == nil {
		atomic.AddInt32(&data.Error, 1)
	}
}

func (r *Collector) Juge(v bool, name ...interface{}) {
	if data, err := r.getExecution(name...); err == nil {
		if v {
			atomic.AddInt32(&data.Success, 1)
		} else {
			atomic.AddInt32(&data.Failed, 1)
		}
	}
}

func (r *Collector) Get() map[string]interface{} {
	return r.data.GetAllAndClear()
}
func (r *Collector) GetByName(name string) interface{} {
	v := r.data.GetAndDel(name)
	if v == nil {
		return struct{}{}
	}
	return v
}
