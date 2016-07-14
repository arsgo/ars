package httpserver

import (
	"sync/atomic"
	"time"
)

//HTTPServerSnap http server快照信息，用于统计执行次数及耗时
type HTTPServerSnap struct {
	Elapsed int64 `json:"elapsed"`
	Min     int64 `json:"min"`
	Max     int64 `json:"max"`
	Count   int64 `json:"count"`
	Average int64 `json:"average"`
}

//Add 添加服务执行时长
func (s *HTTPServerSnap) Add(start time.Time) {
	end := time.Now()
	exp := end.Sub(start).Nanoseconds() / 1000 / 1000
	ce := atomic.AddInt64(&s.Elapsed, exp)
	cc := atomic.AddInt64(&s.Count, 1)
	atomic.SwapInt64(&s.Average, ce/cc)
	if exp < s.Min || s.Min == 0 {
		atomic.SwapInt64(&s.Min, exp)
	}
	if exp > s.Max {
		atomic.SwapInt64(&s.Max, exp)
	}

}
