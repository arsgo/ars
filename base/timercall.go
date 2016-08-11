package base

import "time"

type TimerCall struct {
	msg        chan []interface{}
	timeout    time.Duration
	firstDelay time.Duration
	callback   func(...interface{})
}

func NewTimerCall(timeout time.Duration, firstDelay time.Duration, callback func(...interface{})) *TimerCall {
	tc := &TimerCall{timeout: timeout, callback: callback, firstDelay: firstDelay}
	tc.msg = make(chan []interface{}, 1)
	go tc.call()
	return tc
}
func (t *TimerCall) Push(p ...interface{}) {
	select {
	case t.msg <- p:
	default:
	}
}
func (t *TimerCall) call() {
	time.Sleep(t.firstDelay)
	select {
	case v := <-t.msg:
		t.callback(v...)
	}

	tk := time.NewTicker(t.timeout)
	for {
		select {
		case <-tk.C:
			{
				select {
				case v, ok := <-t.msg:
					if ok {
						t.callback(v...)
					}

				default:
				}
			}
		}
	}
}
