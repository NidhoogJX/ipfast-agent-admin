package cronscheduler

import (
	"github.com/robfig/cron/v3"
)

type JobFunc func()

type Scheduler struct {
	cron *cron.Cron
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithSeconds()),
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) AddJob(spec string, cmd JobFunc) (cron.EntryID, error) {
	return s.cron.AddFunc(spec, cmd)
}
