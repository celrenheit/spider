package schedule

import (
	"time"

	"github.com/gorhill/cronexpr"
)

type CronSchedule struct {
	Expression *cronexpr.Expression
}

// Cron returns a CronSchedule using the cron expression giving as parameter of the function.
func Cron(expression string) CronSchedule {
	expr := cronexpr.MustParse(expression)
	return CronSchedule{
		Expression: expr,
	}
}

func (c CronSchedule) Next(current time.Time) time.Time {
	return c.Expression.Next(current)
}
