package schedulers

import (
	"sort"
	"time"

	"github.com/celrenheit/spider"
	"github.com/celrenheit/spider/spiderutils"
)

type InMemory struct {
	entries Entries
}

func NewInMemory() *InMemory {
	return &InMemory{}
}

type Entry struct {
	Spider   spider.Spider
	Schedule spider.Schedule
	Ctx      *spider.Context
	Next     time.Time
}

type Entries []*Entry

func (e Entries) Len() int      { return len(e) }
func (e Entries) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e Entries) Less(i, j int) bool {
	if e[i].Next.IsZero() {
		return false
	}
	if e[j].Next.IsZero() {
		return true
	}
	return e[i].Next.Before(e[j].Next)
}

func (in *InMemory) Add(sched spider.Schedule, spider spider.Spider) {
	in.AddWithCtx(sched, spider, nil)
}

func (in *InMemory) AddWithCtx(sched spider.Schedule, spider spider.Spider, ctx *spider.Context) {
	entry := &Entry{
		Spider:   spider,
		Schedule: sched,
		Ctx:      ctx,
	}
	in.entries = append(in.entries, entry)
}

func (in *InMemory) AddFunc(sched spider.Schedule, url string, fn func(*spider.Context) error) {
	s := spiderutils.NewGETSpider(url, fn)
	in.AddWithCtx(sched, s, nil)
}

func (in *InMemory) Start() {
	if len(in.entries) == 0 {
		panic("No spider registered")
	}
	now := time.Now()
	for _, e := range in.entries {
		e.Next = e.Schedule.Next(now)
	}
	for {
		sort.Sort(in.entries)
		nextRun := in.entries[0].Next

		select {
		case <-time.After(nextRun.Sub(now)):
			for _, e := range in.entries {
				if e.Next != nextRun {
					break
				}
				ctx, _ := e.Spider.Setup(e.Ctx)
				go e.Spider.Spin(ctx)
				e.Next = e.Schedule.Next(nextRun)
			}
		}
	}
}
