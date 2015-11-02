package spider_test

import (
	"testing"
	"time"

	"github.com/celrenheit/spider"
	"github.com/celrenheit/spider/schedule"
)

func TestInMemory(t *testing.T) {
	ran := false
	testSpider := spider.Get("http://google.com", func(ctx *spider.Context) error {
		ran = true
		return nil
	})
	sched := spider.NewScheduler()
	sched.Add(schedule.Every(1*time.Second), testSpider)
	sched.Start()
	dur := 1*time.Second + 500*time.Millisecond

	select {
	case <-time.After(dur):
		if !ran {
			t.Error("spider not ran")
		}
	}
	sched.Stop()
}

func TestNotRanWhenStopped(t *testing.T) {
	ran := false
	testSpider := spider.Get("http://google.com", func(ctx *spider.Context) error {
		ran = true
		return nil
	})
	dur := 1*time.Second + 100*time.Millisecond
	stopCh := make(chan struct{})

	sched := spider.NewScheduler()
	sched.Add(schedule.Every(1*time.Second), testSpider)
	sched.Start()

	go func() {
		sched.Stop()
		stopCh <- struct{}{}
	}()

	select {
	case <-time.After(dur):
		t.Error("Should not wait to much")
	case <-stopCh:
		if ran {
			t.Error("Spider ran but should not run when stopped")
		}
	}
}
