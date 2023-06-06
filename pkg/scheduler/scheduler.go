package scheduler

import (
	"time"
)

type Scheduler struct {
	timeTick      time.Duration
	funcTick      func() error
	callbackError func(error)
}

func NewScheduler(timeTick time.Duration, f func() error, callbackError func(error)) *Scheduler {
	return &Scheduler{
		funcTick:      f,
		timeTick:      timeTick,
		callbackError: callbackError,
	}
}

func (s *Scheduler) Run() {
	go func(scheduler *Scheduler) {
		for {
			if e := scheduler.funcTick(); e != nil {
				s.callbackError(e)
			}
			time.Sleep(s.timeTick)
		}
	}(s)
}
