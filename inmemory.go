package spider

import (
	"sort"
	"time"
)

// InMemory is the default scheduler
type InMemory struct {
	entries Entries
	addCh   chan *Entry
	stopCh  chan struct{}
	running bool
}

// NewScheduler returns a new InMemory scheduler
func NewScheduler() *InMemory {
	return &InMemory{
		addCh:   make(chan *Entry),
		stopCh:  make(chan struct{}),
		entries: nil,
	}
}

// Entry groups a spider, its root context, a Schedule and the Next time the spider must be launched
type Entry struct {
	Spider   Spider
	Schedule Schedule
	Ctx      *Context
	Next     time.Time
}

// Entries is a collection of Entry.
// Sortable by time.
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

// Add adds a spider using a nil root Context
func (in *InMemory) Add(sched Schedule, spider Spider) {
	in.AddWithCtx(sched, spider, nil)
}

// AddWithCtx adds a spider with a root Context passed in the arguments
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

// AddFunc allows to add a spider using an url and a closure.
// It is by default using the GET HTTP method.
func (in *InMemory) AddFunc(sched Schedule, url string, fn func(*Context) error) {
	s := Get(url, fn)
	in.AddWithCtx(sched, s, nil)
}

// Start launch the scheduler.
// It will run in its own goroutine.
// Your code will continue to be execute after calling this function.
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
		case now = <-time.After(nextRun.Sub(now)):
			for _, e := range in.entries {
				if e.Next != nextRun {
					break
				}
				go in.runEntry(e)
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

func (in *InMemory) runEntry(e *Entry) {
	ctx, _ := e.Spider.Setup(e.Ctx)
	e.Spider.Spin(ctx)
}

// Stop the scheduler.
// Should be called after Start.
func (in *InMemory) Stop() {
	in.stopCh <- struct{}{}
	in.running = false
}

// Standard Scheduler
var stdSched = NewScheduler()

// Add adds a spider to the standard scheduler
func Add(sched Schedule, spider Spider) {
	stdSched.Add(sched, spider)
}

// AddFunc allows to add a spider to the standard scheduler using an url and a closure.
func AddFunc(sched Schedule, url string, fn func(*Context) error) {
	stdSched.AddFunc(sched, url, fn)
}

// Start starts the standard scheduler
func Start() {
	stdSched.Start()
}

// Stop stops the standard scheduler
func Stop() {
	stdSched.Stop()
}
