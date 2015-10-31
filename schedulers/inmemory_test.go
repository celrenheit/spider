package schedulers_test

import (
	"testing"
	"time"

	"github.com/celrenheit/spider"
	"github.com/celrenheit/spider/schedule"
	"github.com/celrenheit/spider/schedulers"
	"github.com/celrenheit/spider/spiderutils"
)

func TestInMemory(t *testing.T) {
	ran := false
	testSpider := spiderutils.NewGETSpider("http://google.com", func(ctx *spider.Context) error {
		ran = true
		return nil
	})
	sched := schedulers.NewInMemory()
	sched.Add(schedule.Every(2*time.Second), testSpider)
	go sched.Start()
	dur := 2*time.Second + 200*time.Millisecond

	select {
	case <-time.After(dur):
		if !ran {
			t.Error("spider not ran")
		}
	}
}
