package spiders

import (
	"io"

	"github.com/celrenheit/spider"
	"github.com/celrenheit/spider/spiderutils"
)

var _ spider.Spider = (*SimpleSpider)(nil)

type SimpleSpider struct {
	URL         string
	Method      string
	Body        io.Reader
	spinnerFunc spider.SpinnerFunc
}

// NewSimpleSpider creates a new SimpleSpider with provided arguments.
// This function requires a spider.SpinnerFunc to handle parsing of to response.
func NewSimpleSpider(method, url string, body io.Reader, spinnerFunc spider.SpinnerFunc) *SimpleSpider {
	return &SimpleSpider{
		Method:      method,
		URL:         url,
		Body:        body,
		spinnerFunc: spinnerFunc,
	}
}

func (s *SimpleSpider) Setup(ctx *spider.Context) (*spider.Context, error) {
	return spiderutils.NewHTTPContext(s.Method, s.URL, s.Body)
}

func (s *SimpleSpider) Spin(ctx *spider.Context) error {
	if s.spinnerFunc == nil {
		panic("SimpleSpider: no spinnerFunc function has been defined")
	}
	return s.spinnerFunc(ctx)
}

// Create a SimpleSpider that makes a get request to the specified url.
func NewSimpleGETSpider(url string, spinnerFunc spider.SpinnerFunc) *SimpleSpider {
	return NewSimpleSpider("GET", url, nil, spinnerFunc)
}
