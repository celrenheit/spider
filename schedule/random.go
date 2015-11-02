package schedule

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

// RandomInterval defines a random interval schedule.
type RandomInterval struct {
	Interval   time.Duration
	Randomness float64
}

// EveryRandom takes an interval with an ajustable plus or minus percentage of this interval.
// The plusOrMinus paramter should be between 0 and 1.
// It returns a Schedule.
// For example, EveryRandom(4*time.Second, 0.5) will return a Schedule that can return between 2 and 6 seconds.
func EveryRandom(interval time.Duration, plusOrMinus float64) RandomInterval {

	if interval < time.Second {
		interval = time.Second
	}

	if plusOrMinus > 1 {
		plusOrMinus = 1
	}

	if plusOrMinus < 0 {
		plusOrMinus = 0
	}

	return RandomInterval{
		Interval:   interval,
		Randomness: plusOrMinus,
	}
}

func (r RandomInterval) RandomRange() (min, max float64) {
	min = r.Interval.Seconds() - r.Interval.Seconds()*r.Randomness
	max = r.Interval.Seconds() + r.Interval.Seconds()*r.Randomness
	return min, max
}

func (r RandomInterval) Next(now time.Time) time.Time {
	min, max := r.RandomRange()
	return now.Add(time.Duration(randomInt(min, max)) * time.Second)
}

func randomInt(min, max float64) int64 {
	return rand.Int63n(int64(max-min)) + int64(min)
}
