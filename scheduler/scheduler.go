package scheduler

import "github.com/colinyl/cron"

type Scheduler struct {
	c *cron.Cron
}

var (
	Schd *Scheduler = NewScheduler()
)

func NewScheduler() *Scheduler {
	return &Scheduler{c: cron.New()}
}

func AddTask(trigger string, task *TaskDetail) {
	Schd.c.AddJob(trigger, task)
}
func Start() {
	Schd.c.Start()
}

func Stop() {
	Schd.c.Stop()
	Schd = &Scheduler{c: cron.New()}
}

func (s *Scheduler) AddTask(trigger string, task *TaskDetail) {
	s.c.AddJob(trigger, task)
}
func (s *Scheduler) Start() {
	s.c.Start()
}

func (s *Scheduler) Stop() {
	s.c.Stop()
	s.c = cron.New()
}
