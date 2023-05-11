package scheduler

import "fmt"

type Scheduler struct {
	funcTick func() error
}

func NewScheduler(f func() error) *Scheduler {
	return &Scheduler{
		funcTick: f,
	}
}

func (s *Scheduler) Run() {
	go func(scheduler *Scheduler) {
		for {
			if e := scheduler.funcTick(); e != nil {
				fmt.Println(e)
			}
		}
	}(s)
}
