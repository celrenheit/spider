package spider

import "time"

// Scheduler is an interface defining an interface that scheduler.
// To define your own scheduler you should implement this interface
type Scheduler interface {
	Handle(Spider) SpiderScheduler
	Start() error
}

// BaseSpiderScheduler is an interface that represents the core methods need for a spider scheduler.
type BaseSpiderScheduler interface {
	NextSpin() (time.Duration, bool)
	NextSpinChan() (<-chan struct{}, <-chan struct{})
}

// SpiderScheduler is an interface that allow to specify a schedule for the current spider added to the scheduler
type SpiderScheduler interface {
	// Definition
	Every(time.Duration) SpiderScheduler
	EveryFunc(EveryFunc) SpiderScheduler
	EveryRandom(time.Duration, time.Duration, time.Duration) SpiderScheduler
	From(time.Time) SpiderScheduler
	To(time.Time) SpiderScheduler
	After(time.Duration) SpiderScheduler
	Delay() time.Duration
	Duplicate(int64) SpiderScheduler
	NumGoroutine() int64

	// Base
	BaseSpiderScheduler
}

type EveryFunc func() time.Duration
