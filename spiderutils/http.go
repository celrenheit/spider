package spiderutils

import (
	"io"
	"net/http"

	"github.com/celrenheit/spider"
)

// NewHTTPContext returns a new spider.Context.
//
// It creates a new http.Client and a new http.Request with the provided arguments.
func NewHTTPContext(method, url string, body io.Reader) (*spider.Context, error) {
	ctx := spider.NewContext()
	// Setup client
	if _, err := ctx.NewClient(); err != nil {
		return ctx, err
	}
	// Request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return ctx, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	ctx.SetRequest(req)
	return ctx, nil
}
