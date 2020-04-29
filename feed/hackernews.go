package feed

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Scraper interface {
	GetInitialFeed() chan Feed
	GetNextFeed() chan Feed
	GetPrevFeed() chan Feed
	GetFeedName() string
}

// todo: add sitestr
type Feed struct {
	Title string
	Link  string
}

type scraper struct {
	feedBase   string
	newestPath string
	nextPath   string
	prevPaths  []string
	pageIndex  int
}

func NewHackerNewsScraper() Scraper {
	return &scraper{
		feedBase:   "https://news.ycombinator.com/",
		newestPath: "newest",
	}
}

func (s *scraper) GetFeedName() string {
	return fmt.Sprintf("  Hacker News Feed [%d]  ", s.pageIndex)
}

func (s *scraper) GetInitialFeed() chan Feed {
	s.pageIndex = 1
	s.prevPaths = make([]string, 0)
	s.prevPaths = append(s.prevPaths, s.newestPath)
	return s.scrapeFeed(fmt.Sprintf("%s%s", s.feedBase, s.newestPath))
}

func (s *scraper) GetNextFeed() chan Feed {
	s.pageIndex += 1
	if s.nextPath != "" && len(s.prevPaths) < s.pageIndex {
		s.prevPaths = append(s.prevPaths, s.nextPath)
	}
	return s.scrapeFeed(fmt.Sprintf("%s%s", s.feedBase, s.nextPath))
}

func (s *scraper) GetPrevFeed() chan Feed {
	if s.pageIndex > 1 {
		s.pageIndex -= 1
	}
	return s.scrapeFeed(fmt.Sprintf("%s%s", s.feedBase, s.prevPaths[s.pageIndex-1]))
}

func (s *scraper) scrapeFeed(url string) chan Feed {
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
			s.nextPath = url
		})

		collector.Visit(url)
	}()

	return ch
}
