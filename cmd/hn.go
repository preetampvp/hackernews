package main

import (
	_ "github.com/pkg/browser"
	"github.com/preetampvp/hackernews/feed"
	"github.com/preetampvp/hackernews/ui"
)

func main() {
	feed := feed.NewHackerNewsScraper()
	ui := ui.NewFeedViewer(feed)
	ui.Show()

	// _ = browser.OpenURL("https://www.nytimes.com/2013/02/17/fashion/creating-hipsturbia-in-the-suburbs-of-new-york.html")
}
