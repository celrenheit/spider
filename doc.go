// Package spider offers a way to scrape websites.

// Installation:
//
//    go get -u github.com/celrenheit/spider
//
//
// Usage of this package is around the usage of spiders and passing contexts.
//
//
//    ctx, err := spider.Setup(nil)
//    err := spider.Spin(ctx)
//
// If you have many spider you can make use of a scheduler. This package provides a basic scheduler.
//
//    scheduler := spider.NewScheduler()
//
//    scheduler.Add(schedule.Every(20 * time.Second), spider1)
//
//    scheduler.Add(schedule.Every(20 * time.Second),spider2)
//
//    scheduler.Start()
//
// This will launch 2 spiders every 20 seconds for the first and every 10 seconds for the second.
//
//
// You can create you own spider by implementing the Spider interface
//
//
//    package main
//
//    import (
//    	"fmt"
//
//    	"github.com/celrenheit/spider"
//    )
//
//    func main() {
//    	wikiSpider := &WikipediaHTMLSpider{
//    		Title: "Albert Einstein",
//    	}
//    	ctx, _ := wikiSpider.Setup(nil)
//    	wikiSpider.Spin(ctx)
//    }
//
//    type WikipediaHTMLSpider struct {
//    	Title string
//    }
//
//    func (w *WikipediaHTMLSpider) Setup(ctx *spider.Context) (*spider.Context, error) {
//    	url := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", w.Title)
//    	return spider.NewHTTPContext("GET", url, nil)
//    }
//
//    func (w *WikipediaHTMLSpider) Spin(ctx *spider.Context) error {
//    	if _, err := ctx.DoRequest(); err != nil {
//    		return err
//    	}
//
//    	html, _ := ctx.HTMLParser()
//    	summary := html.Find("#mw-content-text p").First().Text()
//
//    	fmt.Println(summary)
//    	return nil
//    }
//
//
package spider
