package schedule

import (
	"testing"
	"time"
)

func TestBelowOneSecond(t *testing.T) {
	duration := 200 * time.Nanosecond
	s := Every(duration)
	if s.Interval != time.Second {
		t.Errorf("Expecting %s to be one second:", duration)
	}
}

func TestRoundedToTheSecond(t *testing.T) {
	duration := 2*time.Second + 200*time.Nanosecond
	s := Every(duration)
	if s.Interval != 2*time.Second {
		t.Errorf("Expecting %s to be two second:", duration)
	}
}

func TestNext(t *testing.T) {
	duration := 4 * time.Second
	s := Every(duration)
	now := time.Now().Round(1 * time.Second)
	if s.Next(now) != now.Add(4*time.Second) {
		t.Errorf("Expected: %s but got %s", now, s.Next(now))
	}
}
