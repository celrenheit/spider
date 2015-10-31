package schedule

import "time"

type ConstantSchedule struct {
	Interval time.Duration
}

func Every(duration time.Duration) ConstantSchedule {
	if duration < time.Second {
		duration = time.Second
	}
	duration = duration - time.Duration(duration.Nanoseconds())%time.Second
	return ConstantSchedule{
		Interval: duration,
	}
}

func (c ConstantSchedule) Next(current time.Time) time.Time {
	return current.Add(c.Interval).Round(1 * time.Second)
}
