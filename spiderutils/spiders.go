package spiderutils

import (
	"io"

	"github.com/celrenheit/spider"
)

type spinFunc func(*spider.Context) error

type spiderFunc struct {
	method string
	url    string
	body   io.Reader
	fn     spinFunc
}

func (s *spiderFunc) Setup(parent *spider.Context) (*spider.Context, error) {
	return NewHTTPContext(s.method, s.url, s.body)
}
func (s *spiderFunc) Spin(ctx *spider.Context) error { return s.fn(ctx) }

func NewHTTPSpider(method, url string, body io.Reader, fn spinFunc) *spiderFunc {
	return &spiderFunc{
		method: method,
		url:    url,
		body:   body,
		fn:     fn,
	}
}

func NewGETSpider(url string, fn spinFunc) *spiderFunc {
	return NewHTTPSpider("GET", url, nil, fn)
}

func NewPOSTSpider(url string, body io.Reader, fn spinFunc) *spiderFunc {
	return NewHTTPSpider("POST", url, body, fn)
}

func NewPUTSpider(url string, body io.Reader, fn spinFunc) *spiderFunc {
	return NewHTTPSpider("PUT", url, body, fn)
}

func NewDELETESpider(url string, fn spinFunc) *spiderFunc {
	return NewHTTPSpider("DELETE", url, nil, fn)
}
