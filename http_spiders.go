package spider

import "io"

type spinFunc func(*Context) error

type spiderFunc struct {
	method string
	url    string
	body   io.Reader
	fn     spinFunc
}

func (s *spiderFunc) Setup(parent *Context) (*Context, error) {
	return NewHTTPContext(s.method, s.url, s.body)
}
func (s *spiderFunc) Spin(ctx *Context) error { return s.fn(ctx) }

// NewHTTPSpider creates a new spider according to the http method, url and body.
// The last argument is a closure for doing the actual work
func NewHTTPSpider(method, url string, body io.Reader, fn spinFunc) *spiderFunc {
	return &spiderFunc{
		method: method,
		url:    url,
		body:   body,
		fn:     fn,
	}
}

// Get returns a new GET HTTP Spider.
func Get(url string, fn spinFunc) *spiderFunc {
	return NewHTTPSpider("GET", url, nil, fn)
}

// Post returns a new POST HTTP Spider.
func Post(url string, body io.Reader, fn spinFunc) *spiderFunc {
	return NewHTTPSpider("POST", url, body, fn)
}

// Put returns a new PUT HTTP Spider.
func Put(url string, body io.Reader, fn spinFunc) *spiderFunc {
	return NewHTTPSpider("PUT", url, body, fn)
}

// Delete returns a new DELETE HTTP Spider.
func Delete(url string, fn spinFunc) *spiderFunc {
	return NewHTTPSpider("DELETE", url, nil, fn)
}
