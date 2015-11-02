package spider

import "time"

// Spider is an interface with two methods.
// It is the primary element of the package
type Spider interface {
	Setup(*Context) (*Context, error)
	Spin(*Context) error
}

// Schedule is an interface with only a Next method.
// Next will return the next time it should run given the current time as a parameter.
type Schedule interface {
	Next(time.Time) time.Time
}
