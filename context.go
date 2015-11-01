package spider

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"github.com/cenkalti/backoff"
	"golang.org/x/net/publicsuffix"
)

var (
	ErrNoClient  = errors.New("No request has been set")
	ErrNoRequest = errors.New("No request has been set")
)

// Context is the element that can be shared accross different spiders.
// It contains an HTTP Client and an HTTP Request.
// Context can execute an HTTP Request.
type Context struct {
	Client   *http.Client
	response *http.Response
	request  *http.Request
	Parent   *Context
	Children []*Context
	store    *store
}

// NewContext returns a new Context.
func NewContext() *Context {
	return &Context{
		store:    NewKVStore(),
		Children: make([]*Context, 0),
	}
}

// NewHTTPContext returns a new Context.
//
// It creates a new http.Client and a new http.Request with the provided arguments.
func NewHTTPContext(method, url string, body io.Reader) (*Context, error) {
	ctx := NewContext()
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

// HTMLParser returns an HTML parser.
//
// It uses PuerkitoBio's awesome goquery package.
// It can be found an this url: https://github.com/PuerkitoBio/goquery.
func (c *Context) HTMLParser() (*goquery.Document, error) {
	res := c.Response()
	defer res.Body.Close()
	return goquery.NewDocumentFromReader(res.Body)
}

// JSONParser returns a JSON parser.
//
// It uses Bitly's go-simplejson package which can be found in: https://github.com/bitly/go-simplejson
func (c *Context) JSONParser() (*simplejson.Json, error) {
	res := c.Response()
	defer res.Body.Close()
	return simplejson.NewFromReader(res.Body)
}

// RAWContent returns the raw data of the reponse's body
func (c *Context) RAWContent() ([]byte, error) {
	res := c.Response()
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// Response returns an http.Response
func (c *Context) Response() *http.Response {
	return c.response
}

// SetResponse set an http.Response
func (c *Context) SetResponse(res *http.Response) {
	c.response = res
}

// SetRequest set an http.Request
func (c *Context) SetRequest(req *http.Request) {
	c.request = req
}

// Request returns an http.Response
func (c *Context) Request() *http.Request {
	return c.request
}

// Cookies return a list of cookies for the given request URL
func (c *Context) Cookies() []*http.Cookie {
	return c.Client.Jar.Cookies(c.Request().URL)
}

// ResetCookies create a new cookie jar.
//
// Note: All the cookies previously will be deleted.
func (c *Context) ResetCookies() error {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return err
	}
	c.Client.Jar = jar
	return nil
}

// ResetClient create a new http.Client and replace the existing one if there is one.
func (c *Context) ResetClient() (*http.Client, error) {
	client, err := c.NewClient()
	c.Client = client
	return client, err
}

// NewClient create a new http.Client
func (c *Context) NewClient() (*http.Client, error) {
	jar, err := c.NewCookieJar()
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar:       jar,
		Transport: http.DefaultTransport,
	}
	c.Client = client
	return client, nil
}

// NewCookieJar create a new *cookiejar.Jar
func (c *Context) NewCookieJar() (*cookiejar.Jar, error) {
	return cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
}

// DoRequest makes an http request using the http.Client and http.Request associated with this context.
//
// This will store the response in this context. To access the response you should do:
// 		ctx.Response() // to get the http.Response
func (c *Context) DoRequest() (*http.Response, error) {
	client := c.Client
	if client == nil {
		client = http.DefaultClient
	}
	if c.Request() == nil {
		return nil, ErrNoRequest
	}
	res, err := client.Do(c.Request())
	if err == nil {
		c.SetResponse(res)
	}
	return res, err
}

// DoRequestWithExponentialBackOff makes an http request using the http.Client and http.Request associated with this context.
// You can pass a condition and a BackOff configuration. See https://github.com/cenkalti/backoff to know more about backoff.
// If no BackOff is provided it will use the default exponential BackOff configuration.
// See also ErrorIfStatusCodeIsNot function that provides a basic condition based on status code.
func (c *Context) DoRequestWithExponentialBackOff(condition BackoffCondition, b backoff.BackOff) (*http.Response, error) {
	if b == nil {
		b = backoff.NewExponentialBackOff()
	}
	err := backoff.RetryNotify(
		func() error {
			var err error
			res, err := c.DoRequest()
			if err != nil {
				return err
			}
			return condition(res)
		},
		b,
		func(err error, wait time.Duration) {
			fmt.Println("˙\nBackoff:Waiting: ", wait)
			fmt.Println(err)
			fmt.Println("")
		})
	return c.Response(), err
}

// Set a value to this context
func (c *Context) Set(key string, value interface{}) {
	c.store.set(key, value)
}

// Get a value from this context
func (c *Context) Get(key string) interface{} {
	return c.store.get(key)
}

// ExtendWithRequest return a new Context child to the provided context associated with the provided http.Request.
func (c *Context) ExtendWithRequest(ctx Context, r *http.Request) *Context {
	newCtx := NewContext()
	newCtx.SetRequest(r)
	newCtx.Parent = c
	return newCtx
}

// Set a parent context to the current context.
// It will also add the current context to the list of children of the parent context.
func (c *Context) SetParent(parent *Context) {
	c.Parent = parent
	parent.Children = append(parent.Children, c)
}

// NewKVStore returns a new store.
func NewKVStore() *store {
	return &store{
		kv: make(map[string]interface{}),
	}
}

type store struct {
	sync.RWMutex
	kv map[string]interface{}
}

func (s *store) set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()
	s.kv[key] = value
}

func (s *store) get(key string) interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.kv[key]
}

type BackoffCondition func(*http.Response) error

func ErrorIfStatusCodeIsNot(status int) BackoffCondition {
	return func(res *http.Response) error {
		if res.StatusCode != status {
			return errors.New("Request failed: " + res.Status)
		}
		return nil
	}
}
