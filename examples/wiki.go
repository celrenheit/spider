package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/celrenheit/spider"
	"github.com/celrenheit/spider/schedulers"
	"github.com/celrenheit/spider/spiderutils"
)

var (
	// Ensure WikipediaHTMLSpider implements spider.Spider interface
	_ spider.Spider = (*WikipediaHTMLSpider)(nil)
	// Ensure WikipediaJSONSpider implements spider.Spider interface
	_ spider.Spider = (*WikipediaJSONSpider)(nil)
)

func main() {
	wikiHTMLSpider := &WikipediaHTMLSpider{"Albert Einstein"}
	wikiJSONSpider := &WikipediaJSONSpider{"Lionel Messi"}

	s := schedulers.NewBasicScheduler()
	s.Handle(wikiHTMLSpider).Every(30 * time.Second)
	s.Handle(wikiJSONSpider).Every(20 * time.Second)
	s.Start()
}

type WikipediaHTMLSpider struct {
	Title string
}

func (w *WikipediaHTMLSpider) Setup(ctx *spider.Context) (*spider.Context, error) {
	url := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", w.Title)
	return spiderutils.NewHTTPContext("GET", url, nil)
}

func (w *WikipediaHTMLSpider) Spin(ctx *spider.Context) error {
	if _, err := ctx.DoRequest(); err != nil {
		return err
	}

	html, _ := ctx.HTMLParser()
	summary := html.Find("#mw-content-text p").First().Text()

	fmt.Println(summary)
	return nil
}

type WikipediaJSONSpider struct {
	Title string
}

func (w *WikipediaJSONSpider) Setup(ctx *spider.Context) (*spider.Context, error) {
	params := url.Values{}
	params.Add("titles", w.Title)
	url := fmt.Sprintf("http://en.wikipedia.org/w/api.php?format=json&action=query&prop=extracts&exintro=&explaintext=&%s", params.Encode())
	return spiderutils.NewHTTPContext("GET", url, nil)
}

func (w *WikipediaJSONSpider) Spin(ctx *spider.Context) error {
	if _, err := ctx.DoRequest(); err != nil {
		return err
	}
	jsonparser, _ := ctx.JSONParser()
	pages, err := jsonparser.GetPath("query", "pages").Map()
	if err != nil {
		return err
	}
	for _, p := range pages {
		fmt.Println(p.(map[string]interface{})["extract"])
	}
	return nil
}
