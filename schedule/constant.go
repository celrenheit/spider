package schedule

import "time"

type ConstantSchedule struct {
	Interval time.Duration
}

// Every returns a ConstantSchedule that runs every duration given as parameter.
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
	return current.Add(c.Interval - time.Duration(current.Nanosecond())*time.Nanosecond)
}
