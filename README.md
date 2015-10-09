# Spider [![Build Status](https://travis-ci.org/celrenheit/spider.svg?branch=master)](https://travis-ci.org/celrenheit/spider) [![GoDoc](https://godoc.org/github.com/celrenheit/spider?status.svg)](https://godoc.org/github.com/celrenheit/spider) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Spider package provides some simple spider and scheduler interfaces for scraping and parsing HTML and JSON pages.

# Installation

```shell
$ go get -u github.com/celrenheit/spider
```

# Documentation

The documentation is hosted on [GoDoc](https://godoc.org/github.com/celrenheit/spider).

# Usage

```go
package main

import (
	"fmt"
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
	scheduler.Start()
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
	htmlparser, _ := ctx.HTMLParser()
	// Get the first paragraph of the wikipedia page
	summary := htmlparser.Find("#mw-content-text p").First().Text()

	fmt.Println(summary)
	return nil
}
```

# Contributing

Contributions are welcome ! Feel free to submit a pull request.
You can improve documentation and examples to start.
You can also provides spiders and better schedulers.

# License

[MIT License](https://github.com/celrenheit/spider/blob/master/LICENSE)
