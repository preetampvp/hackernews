package main

import (
	"fmt"

	_ "github.com/pkg/browser"
	"github.com/preetampvp/hackernews/feed"
)

func main() {
	feedBase := "https://news.ycombinator.com/"
	feed := feed.NewScraper()

	for f := range feed.ScrapeFeed(fmt.Sprintf("%snewest", feedBase)) {
		// go browser.OpenURL(f.Link)
		fmt.Println(f.Title, f.Link)
	}
}
