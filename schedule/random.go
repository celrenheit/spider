package schedule

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

type RandomInterval struct {
	Interval   time.Duration
	Randomness float64
}

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
