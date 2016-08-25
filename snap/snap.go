package snap

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/arsgo/lib4go/concurrent"
)

//Snap 快照信息
type Snap struct {
	data        *concurrent.ConcurrentMap
	last        time.Time
	timeout     time.Duration
	lock        sync.Mutex
	callback    []func(map[string]interface{})
	resetTicker chan time.Duration
	isStart     int32
}

var snap *Snap

func init() {
	snap = &Snap{isStart: 1, timeout: 300}
	snap.data = concurrent.NewConcurrentMap()
	snap.resetTicker = make(chan time.Duration, 1)
}

//Bind 绑定回调函数
func Bind(seconds time.Duration, callback func(map[string]interface{})) {
	if atomic.CompareAndSwapInt32(&snap.isStart, 1, 0) {
		go call()
	}
	snap.timeout = seconds
	snap.lock.Lock()
	defer snap.lock.Unlock()
	snap.callback = append(snap.callback, callback)
}

//ResetTicker 重置ticker
func ResetTicker(seconds time.Duration) {
	if snap.timeout != seconds && snap.isStart == 0 {
		snap.timeout = seconds
		snap.resetTicker <- seconds
	}
}

func call() {
	tk := time.NewTicker(snap.timeout)
	for {
		select {
		case <-tk.C:
			callbackNow()
		case t := <-snap.resetTicker:
			tk.Stop()
			tk = time.NewTicker(t)
		}
	}
}
func callbackNow() {
	snap.lock.Lock()
	defer snap.lock.Unlock()
	data := GetData()
	for _, v := range snap.callback {
		go v(data)
	}
}
func GetData() map[string]interface{} {
	return snap.data.GetAllAndClear()
}

//Append 将指定数据添加到快照集合
func Append(name string, data interface{}) {
	snap.data.Set(name, data)
}
