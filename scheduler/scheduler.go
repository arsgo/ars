package scheduler

import "github.com/colinyl/cron"

type Scheduler struct {
	c *cron.Cron
}

var (
	Schd *Scheduler = &Scheduler{c: cron.New()}
)

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
