package base

import "sync"

type Sync struct {
	sync    *sync.WaitGroup
	lk      *sync.Mutex
	actions map[string]int
}

func NewSync(wait int) Sync {
	sn := Sync{}
	sn.lk = &sync.Mutex{}
	sn.sync = &sync.WaitGroup{}
	sn.actions = make(map[string]int)
	sn.sync.Add(wait)
	return sn
}
func (s Sync) Wait() {
	s.sync.Wait()
}
func (s Sync) WaitAndAdd(delta int) {
	s.sync.Wait()
	s.sync.Add(delta)
}

func (s Sync) AddStep(delta int) {
	s.sync.Add(delta)
}
func (s Sync) Done(action string) {
	s.lk.Lock()
	defer s.lk.Unlock()
	if _, ok := s.actions[action]; !ok {
		s.actions[action] = 1
		s.sync.Done()
	}
}
