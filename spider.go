package spider

import "time"

// Spider is an interface with two methods.
// It is the primary element of the package
type Spider interface {
	Setup(*Context) (*Context, error)
	Spin(*Context) error
}

type SpinnerFunc func(ctx *Context) error

func (s SpinnerFunc) Spin(ctx *Context) error {
	return s(ctx)
}

type Schedule interface {
	Next(time.Time) time.Time
}
