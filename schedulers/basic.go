package schedulers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/celrenheit/spider"
)

var _ spider.Scheduler = (*BasicScheduler)(nil)

// BasicScheduler is a basic Scheduler that implements the Scheduler interface.
type BasicScheduler struct {
	sync.RWMutex
	schedules map[spider.Spider]spider.SpiderScheduler
	contexts  map[spider.Spider]*spider.Context
}

// NewBasicScheduler returns a new BasicScheduler.
func NewBasicScheduler() *BasicScheduler {
	return &BasicScheduler{
		schedules: make(map[spider.Spider]spider.SpiderScheduler),
		contexts:  make(map[spider.Spider]*spider.Context),
	}
}

// Handle a spider.Spider and returns a new BasicSpiderScheduler associated with this spider.
func (bs *BasicScheduler) Handle(s spider.Spider) spider.SpiderScheduler {
	schedule := NewBasicSpiderScheduler()
	bs.addSpider(s, schedule, nil)
	return schedule
}

// addSpider adds a spider.Spider to the list of spiders.
func (bs *BasicScheduler) addSpider(s spider.Spider, schedule spider.SpiderScheduler, baseCtx *spider.Context) (spider.Spider, error) {
	bs.Lock()
	defer bs.Unlock()
	bs.schedules[s] = schedule
	ctx, err := s.Setup(baseCtx)
	if err != nil {
		return nil, err
	}
	bs.contexts[s] = ctx
	return s, nil
}

// Start starts the Scheduler.
// It will dispatch each spiders in their own watchers for each duplicated (goroutines).
func (bs *BasicScheduler) Start() error {
	errChan := make(chan error)
	doneChan := make(chan struct{})

	var wg sync.WaitGroup
	for s, schedule := range bs.schedules {
		wg.Add(1)
		go func(s spider.Spider, schedule spider.SpiderScheduler, errChan chan error, doneChan chan struct{}) {
			defer wg.Done()
			for i := 0; i < int(schedule.NumGoroutine()); i++ {
				go bs.watchSchedule(s, schedule, errChan, doneChan)
				time.Sleep(schedule.Delay())
			}
		}(s, schedule, errChan, doneChan)
	}
	wg.Wait()

	// Capture and return error
	// Or Count until all spiders have finished
	var doneCount int64 = 0
	total := bs.totalRootSpiders()
	for {
		select {
		case err := <-errChan:
			log.Println(err)
			return err
		case <-doneChan:
			doneCount++
			log.Println(doneCount, "spiders are done. Still", total-doneCount)
			if doneCount == total {
				return nil
			}
		}
	}
}

func (bs *BasicScheduler) totalRootSpiders() int64 {
	var count int64 = 0
	for _, schedule := range bs.schedules {
		count += schedule.NumGoroutine()
	}
	return count
}

func (bs *BasicScheduler) spin(s spider.Spider, baseCtx *spider.Context, errChan chan error) {
	if err := s.Spin(baseCtx); err != nil {
		errChan <- err
	}
}

func (bs *BasicScheduler) watchSchedule(s spider.Spider, schedule spider.SpiderScheduler, errChan chan error, doneChan chan struct{}) {
	bs.RLock()
	c := bs.contexts[s]
	bs.RUnlock()
	spinChan, doneScheduleChan := schedule.NextSpinChan()
	errScheduleChan := make(chan error)
	for {
		select {
		case <-spinChan:
			go bs.spin(s, c, errScheduleChan)
		case err := <-errScheduleChan:
			errChan <- err
			return
		case <-doneScheduleChan:
			fmt.Println("Done !")
			doneChan <- struct{}{}
			return
		}
	}
}
