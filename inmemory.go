package spider

import (
	"sort"
	"time"
)

type InMemory struct {
	entries Entries
	addCh   chan *Entry
	stopCh  chan struct{}
	running bool
}

func NewScheduler() *InMemory {
	return &InMemory{
		addCh:   make(chan *Entry),
		stopCh:  make(chan struct{}),
		entries: nil,
	}
}

type Entry struct {
	Spider   Spider
	Schedule Schedule
	Ctx      *Context
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

func (in *InMemory) Add(sched Schedule, spider Spider) {
	in.AddWithCtx(sched, spider, nil)
}

func (in *InMemory) AddWithCtx(sched Schedule, spider Spider, ctx *Context) {
	entry := &Entry{
		Spider:   spider,
		Schedule: sched,
		Ctx:      ctx,
	}
	if !in.running {
		in.entries = append(in.entries, entry)
		return
	}
	in.addCh <- entry
}

func (in *InMemory) AddFunc(sched Schedule, url string, fn func(*Context) error) {
	s := Get(url, fn)
	in.AddWithCtx(sched, s, nil)
}

func (in *InMemory) Start() {
	in.running = true
	go in.start()
}

func (in *InMemory) start() {
	now := time.Now().Local()
	for _, e := range in.entries {
		e.Next = e.Schedule.Next(now)
	}
	for {
		sort.Sort(in.entries)
		var nextRun time.Time

		if len(in.entries) == 0 {
			// Wait 1 day if there is no spiders to run
			nextRun = now.Add(24 * time.Hour)
		} else {
			nextRun = in.entries[0].Next
		}

		select {
		case <-time.After(nextRun.Sub(now)):
			for _, e := range in.entries {
				if e.Next != nextRun {
					break
				}
				go func(e *Entry) {
					ctx, _ := e.Spider.Setup(e.Ctx)
					go e.Spider.Spin(ctx)
				}(e)
				e.Next = e.Schedule.Next(nextRun)
			}
			continue
		case e := <-in.addCh:
			in.entries = append(in.entries, e)
			e.Next = e.Schedule.Next(now)
		case <-in.stopCh:
			return
		}
		now = time.Now().Local()
	}
}

func (in *InMemory) Stop() {
	in.stopCh <- struct{}{}
	in.running = false
}

// Standard Scheduler
var stdSched = NewScheduler()

func Add(sched Schedule, spider Spider) {
	stdSched.Add(sched, spider)
}

func AddFunc(sched Schedule, url string, fn func(*Context) error) {
	stdSched.AddFunc(sched, url, fn)
}

func Start() {
	stdSched.Start()
}

func Stop() {
	stdSched.Stop()
}
