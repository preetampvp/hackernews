package feed

import "github.com/gocolly/colly"

type Feed struct {
	Title string
	Link  string
}

type Scraper interface {
	ScrapeFeed(url string) chan Feed
}

func NewScraper() Scraper {
	return &scraper{}
}

type scraper struct{}

// ScrapeFeed
func (s *scraper) ScrapeFeed(url string) chan Feed {
	ch := make(chan Feed)

	go func() {
		collector := colly.NewCollector()
		defer close(ch)

		collector.OnHTML(".storylink", func(e *colly.HTMLElement) {

			title := e.Text
			url := e.Attr("href")

			ch <- Feed{Title: title, Link: url}
		})

		collector.OnHTML(".morelink", func(e *colly.HTMLElement) {

			url := e.Attr("href")

			ch <- Feed{Title: "--MORE--", Link: url}
		})

		collector.Visit(url)
	}()

	return ch
}
