package schedulers

import (
	"math/rand"
	"time"

	"github.com/celrenheit/spider"
)

var _ spider.SpiderScheduler = (*BasicSpiderScheduler)(nil)

// BasicSpiderScheduler represents a schedule configuration for a spider.
type BasicSpiderScheduler struct {
	last       time.Time
	from       time.Time
	to         time.Time
	delay      time.Duration
	goroutines int64
	everyFunc  spider.EveryFunc
	count      int64
}

// NewBasicSpiderScheduler returns a new BasicSpiderScheduler
func NewBasicSpiderScheduler() *BasicSpiderScheduler {
	return &BasicSpiderScheduler{
		goroutines: 1,
	}
}

// func (b *BasicSpiderScheduler) Every() time.Duration {
// 	return 1 * time.Second
// }

// NextSpin returns when a spider has to spin a web based on the provided schedule.
// The second argument is a boolean that indicates if it there is no more spins.
func (b *BasicSpiderScheduler) NextSpin() (time.Duration, bool) {
	var nextSpin time.Duration = 0 * time.Second
	now := time.Now()
	estimatedTime := now

	if b.everyFunc == nil { // Executed only once
		if b.count == 0 { // Execute once
			return nextSpin, true
		} else { // Stop executing
			return nextSpin, false
		}
	}

	if b.everyFunc != nil {
		// Between boundaries
		nextSpin = b.everyFunc()
		estimatedTime.Add(nextSpin)
	}

	if !b.from.IsZero() && b.from.After(estimatedTime) {
		nextSpin = b.from.Sub(estimatedTime)
	} else if !b.to.IsZero() && b.to.Before(estimatedTime) {
		nextSpin = b.to.Sub(estimatedTime)
		return nextSpin, false
	}

	return nextSpin, true
}

// NextSpinChan returns two channels the first one is a spin channel representing
// when a spider have to spin a web and the second is a done channel
// representing when there is webs to be spun.
func (b *BasicSpiderScheduler) NextSpinChan() (<-chan struct{}, <-chan struct{}) {
	spinChan := make(chan struct{})
	doneChan := make(chan struct{})

	go func() {
		for {
			nextSpinDuration, ok := b.NextSpin()
			if !ok { // Done
				doneChan <- struct{}{}
				return
			}
			<-time.After(nextSpinDuration)
			spinChan <- struct{}{}
			b.count++
		}
	}()

	return spinChan, doneChan
}

// Every sets a fixed duration representing when to spin.
func (b *BasicSpiderScheduler) Every(d time.Duration) spider.SpiderScheduler {
	b.everyFunc = func() time.Duration { return d }
	return b
}

// Every sets a function that returns a duration representing when to spin.
func (b *BasicSpiderScheduler) EveryFunc(fn spider.EveryFunc) spider.SpiderScheduler {
	b.everyFunc = fn
	return b
}

// Every sets a randomized duration representing when to spin.
// The second and third argument represents min and max duration for randomness.
// For exemple:
// 		EveryRandom(2*time.Minute, 15*time.Second, 45*time.Second)
//
// Will represent:
// 2 minutes + a random duration between 15 seconds and 45 seconds
func (b *BasicSpiderScheduler) EveryRandom(d, min, max time.Duration) spider.SpiderScheduler {
	b.everyFunc = func() time.Duration {
		return d + time.Duration(randomInt(min.Nanoseconds(), max.Nanoseconds()))
	}
	return b
}

// From defines when the scheduled spider has to start
func (b *BasicSpiderScheduler) From(from time.Time) spider.SpiderScheduler {
	b.from = from
	return b
}

// To defines when the scheduled spider has to stop
func (b *BasicSpiderScheduler) To(to time.Time) spider.SpiderScheduler {
	b.to = to
	return b
}

// After defines a delay between each goroutines
func (b *BasicSpiderScheduler) After(delay time.Duration) spider.SpiderScheduler {
	b.delay = delay
	return b
}

// Delay returns the delay between each goroutines
func (b *BasicSpiderScheduler) Delay() time.Duration {
	return b.delay
}

// Duplicate sets in how many goroutines the current spider will be duplicated.
func (b *BasicSpiderScheduler) Duplicate(numGoRoutines int64) spider.SpiderScheduler {
	b.goroutines = numGoRoutines
	return b
}

// NumGoroutine returns the number of goroutines.
func (b *BasicSpiderScheduler) NumGoroutine() int64 {
	return b.goroutines
}

func randomInt(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}
