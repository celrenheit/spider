package schedule

import (
	"testing"
	"time"
)

func TestRange(t *testing.T) {
	s := EveryRandom(4*time.Second, 0.5)
	min, max := s.RandomRange()
	if min != 2*time.Second.Seconds() {
		t.Errorf("Expected min: %s but got %s", 2*time.Second, min)
	}
	if max != 6*time.Second.Seconds() {
		t.Errorf("Expected min: %s but got %s", 6*time.Second, max)
	}
}

func TestCorrectnessOfRange(t *testing.T) {
	s := EveryRandom(4*time.Second, 0.5)
	now := time.Now().Round(time.Second)
	for i := 0; i < 100000; i++ {
		next := s.Next(now)
		if next.Sub(now) < 2*time.Second || next.Sub(now) > 6*time.Second {
			t.Errorf("Outside of range. Got: %s", next.Sub(now))
		}
	}
}
