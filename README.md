# Spider [![Build Status](https://travis-ci.org/celrenheit/spider.svg?branch=master)](https://travis-ci.org/celrenheit/spider) [![GoDoc](https://godoc.org/github.com/celrenheit/spider?status.svg)](https://godoc.org/github.com/celrenheit/spider) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

This package provides a simple way, yet extensible, to scrape HTML and JSON pages. It uses spiders around the web scheduled at certain configurable intervals to fetch data.
It is written in [Golang](https://golang.org/) and is [MIT licensed](https://github.com/celrenheit/spider#license).

You can see an example app using this package here: [https://github.com/celrenheit/trending-machine](https://github.com/celrenheit/trending-machine)

# Installation

```shell
$ go get -u github.com/celrenheit/spider
```

# Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/celrenheit/spider"
	"github.com/celrenheit/spider/schedule"
)

// LionelMessiSpider scrape wikipedia's page for LionelMessi
// It is defined below in the init function
var LionelMessiSpider spider.Spider

func main() {
	// Create a new scheduler
	scheduler := spider.NewScheduler()

	// Register the spider to be scheduled every 15 seconds
	scheduler.Add(schedule.Every(15*time.Second), LionelMessiSpider)
	// Alternatively, you can choose a cron schedule
	// This will run every minute of every day
	scheduler.Add(schedule.Cron("* * * * *"), LionelMessiSpider)

	// Start the scheduler
	scheduler.Start()

	// Exit 5 seconds later to let time for the request to be done.
	// Depends on your internet connection
	<-time.After(65 * time.Second)
}

func init() {
	LionelMessiSpider = spider.Get("https://en.wikipedia.org/wiki/Lionel_Messi", func(ctx *spider.Context) error {
		fmt.Println(time.Now())
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
		summary := htmlparser.Find("#mw-content-text > p").First().Text()

		fmt.Println(summary)
		return nil
	})
}
```


In order, to create your own spiders you have to implement the [spider.Spider](https://godoc.org/github.com/celrenheit/spider#Spider) interface.
It has two functions, Setup and Spin.

[Setup](https://godoc.org/github.com/celrenheit/spider#Spider) gets a [Context](https://godoc.org/github.com/celrenheit/spider#Context) and returns a new [Context](https://godoc.org/github.com/celrenheit/spider#Context) with an [error](https://godoc.org/builtin#error) if something wrong happened.
Usually, it is in this function that you create a new [http client](https://golang.org/pkg/net/http/#Client) and [http request](https://golang.org/pkg/net/http/#Request).

[Spin](https://godoc.org/github.com/celrenheit/spider#Spider) gets a [Context](https://godoc.org/github.com/celrenheit/spider#Context) do its work and returns an [error](https://godoc.org/builtin#error) if necessarry. It is in this function that you do your work ([do a request](https://godoc.org/github.com/celrenheit/spider#Context.DoRequest), handle response, parse [HTML](https://godoc.org/github.com/celrenheit/spider#Context.HTMLParser) or [JSON](https://godoc.org/github.com/celrenheit/spider#Context.JSONParser), etc...). It should return an error if something didn't happened correctly.


# Documentation

The documentation is hosted on [GoDoc](https://godoc.org/github.com/celrenheit/spider).


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

# Inspiration

[Dkron](https://github.com/victorcoder/dkron) for the new in memory scheduler (as of 0.3)
