# Spider [![Build Status](https://travis-ci.org/celrenheit/spider.svg?branch=master)](https://travis-ci.org/celrenheit/spider) [![GoDoc](https://godoc.org/github.com/celrenheit/spider?status.svg)](https://godoc.org/github.com/celrenheit/spider) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

This package provides a simple way, yet extensible, to scrape HTML and JSON pages. It uses spiders around the web scheduled at certain configurable intervals to fetch data.
It is written in [Golang](https://golang.org/) and is [MIT licensed](https://github.com/celrenheit/spider#license).

# Installation

```shell
$ go get -u github.com/celrenheit/spider
```

# Documentation

The documentation is hosted on [GoDoc](https://godoc.org/github.com/celrenheit/spider).

# Usage

In order, to create your own spiders you have to implement the [spider.Spider](https://godoc.org/github.com/celrenheit/spider#Spider) interface.
It has two functions, Setup and Spin.

[Setup](https://godoc.org/github.com/celrenheit/spider#Spider) gets a [Context](https://godoc.org/github.com/celrenheit/spider#Context) and returns a new [Context](https://godoc.org/github.com/celrenheit/spider#Context) with an [error](https://godoc.org/builtin#error) if something wrong happened.
Usually, it is in this function that you create a new [http client](https://golang.org/pkg/net/http/#Client) and [http request](https://golang.org/pkg/net/http/#Request).

[Spin](https://godoc.org/github.com/celrenheit/spider#Spider) gets a [Context](https://godoc.org/github.com/celrenheit/spider#Context) do its work and returns an [error](https://godoc.org/builtin#error) if necessarry. It is in this function that you do your work ([do a request](https://godoc.org/github.com/celrenheit/spider#Context.DoRequest), handle response, parse [HTML](https://godoc.org/github.com/celrenheit/spider#Context.HTMLParser) or [JSON](https://godoc.org/github.com/celrenheit/spider#Context.JSONParser), etc...). It should return an error if something didn't happened correctly.

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/celrenheit/spider"
	"github.com/celrenheit/spider/schedulers"
	"github.com/celrenheit/spider/spiderutils"
)

func main() {
	wikiSpider := &WikipediaHTMLSpider{"Albert Einstein"}

	// Create a new scheduler
	scheduler := schedulers.NewBasicScheduler()

	// Register the spider to be scheduled every 45 seconds
	scheduler.Handle(wikiSpider).Every(45 * time.Second)

	// Start the scheduler
	log.Fatal(scheduler.Start())
}

type WikipediaHTMLSpider struct {
	Title string
}

func (w *WikipediaHTMLSpider) Setup(ctx *spider.Context) (*spider.Context, error) {
	// Define the url of the wikipedia page
	url := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", w.Title)
	// Create a context with an http.Client and http.Request
	return spiderutils.NewHTTPContext("GET", url, nil)
}

func (w *WikipediaHTMLSpider) Spin(ctx *spider.Context) error {
	// Execute the request
	if _, err := ctx.DoRequest(); err != nil {
		return err
	}

	// Get goquery's html parser
	htmlparser, err := ctx.HTMLParser()
	if err != nil {
		return err
	}
	// Get the first paragraph of the wikipedia page
	summary := htmlparser.Find("#mw-content-text p").First().Text()

	fmt.Println(summary)
	return nil
}
```

# Examples

```shell
$ cd $GOPATH/src/github.com/celrenheit/spider/examples
$ go run wiki.go
```

# Contributing

Contributions are welcome ! Feel free to submit a pull request.
You can improve documentation and examples to start.
You can also provides spiders and better schedulers.

If you have developed your own spiders or schedulers, I will be pleased to review your code and eventually merge it into the project.

# License

[MIT License](https://github.com/celrenheit/spider/blob/master/LICENSE)
