package base

import (
	"errors"
	"sync"
)

type AvaliableMap struct {
	keys  []string
	data  map[string]int
	index int
	lk    sync.Mutex
}

func NewAvaliableMap(keys []string) *AvaliableMap {
	am := &AvaliableMap{keys: keys, index: -1}
	am.data = make(map[string]int)
	for i, v := range am.keys {
		am.data[v] = i
	}
	return am
}
func (a *AvaliableMap) HasAvaliable() bool {
	a.lk.Lock()
	defer a.lk.Unlock()
	return len(a.keys) > 0
}

func (a *AvaliableMap) Get() (s string, err error) {
	a.lk.Lock()
	defer a.lk.Unlock()
	a.index++
	if len(a.keys) > 0 {
		cindex := a.index % len(a.keys)
		s = a.keys[cindex]
		return
	}
	err = errors.New("not avliable")
	return
}
func (a *AvaliableMap) Remove(key string) {
	a.lk.Lock()
	defer a.lk.Unlock()
	if v, ok := a.data[key]; ok {
		delete(a.data, key)
		first := v
		last := v + 1
		var firstLst []string
		var lastLst []string
		if first >= 0 {
			firstLst = a.keys[:first]
		}
		if last < len(a.keys) {
			lastLst = a.keys[last:]
		}
		a.keys = append(firstLst, lastLst...)
	}
	a.data = make(map[string]int)
	for i, v := range a.keys {
		a.data[v] = i
	}
}
