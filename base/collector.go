package base

import (
	"sync/atomic"
	"time"
)

type CollectorCallBack func(s int, f int, e int)

//Collector 采集器
type Collector struct {
	timer    time.Duration
	callback CollectorCallBack
	success  int32
	failed   int32
	err      int32
}

func NewCollector(callback CollectorCallBack, timer time.Duration) *Collector {
	r := &Collector{callback: callback, timer: timer}
	go r.callNow()
	return r
}
func (r *Collector) Success() {
	atomic.AddInt32(&r.success, 1)
}
func (r *Collector) Failed() {
	atomic.AddInt32(&r.failed, 1)
}
func (r *Collector) Error() {
	atomic.AddInt32(&r.err, 1)
}
func (r *Collector) callNow() {
	timer := time.NewTicker(r.timer)
	for {
		select {
		case <-timer.C:
			r.callback(int(atomic.SwapInt32(&r.success, 0)), int(atomic.SwapInt32(&r.failed, 0)), int(atomic.SwapInt32(&r.err, 0)))
		}
	}
}
