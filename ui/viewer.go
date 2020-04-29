package ui

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/preetampvp/hackernews/feed"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type FeedViewer interface {
	Show()
}

type feedLoader func() chan feed.Feed

// NewUi - description
func NewFeedViewer(scraper feed.Scraper) FeedViewer {
	viewer := &feedViewer{scraper: scraper}
	return viewer
}

type feedViewer struct {
	scraper  feed.Scraper
	feed     []feed.Feed
	grid     *ui.Grid
	feedList *widgets.List
	infoText *widgets.Paragraph
}

// Show - Show ui
func (f *feedViewer) Show() {
	if err := ui.Init(); err != nil {
		fmt.Printf("failed to initialize termui: %v", err)
		return
	}
	defer ui.Close()

	f.initListView()
	f.initGrid()
	f.render()
	f.loadFeed(f.scraper.GetInitialFeed)
	f.initEventsPolling()
}

func (f *feedViewer) shortcutsText() string {
	return "[ Shortcuts   ](fg:white,bg:black) [ Enter ](fg:black)[ Open article ](fg:black,bg:green) " +
		"[ q ](fg:black)[ Quit ](fg:black,bg:green) " +
		"[ j ](fg:black)[ Down ](fg:black,bg:green) " +
		"[ k ](fg:black)[ Up ](fg:black,bg:green) " +
		"[ n ](fg:black)[ Next ](fg:black,bg:green) " +
		"[ p ](fg:black)[ Prev ](fg:black,bg:green) " +
		"[ r ](fg:black)[ Refresh ](fg:black,bg:green) "
}

func (f *feedViewer) loadFeed(loader feedLoader) {
	f.infoText.Text = "loading feed...."
	f.render()
	f.feed = make([]feed.Feed, 0)
	f.feedList.Rows = make([]string, 0)
	for item := range loader() {
		f.feed = append(f.feed, item)
		f.feedList.Rows = append(f.feedList.Rows, item.Title)
	}
	f.feedList.Title = f.scraper.GetFeedName()
	f.infoText.Text = f.shortcutsText()
	f.feedList.SelectedRow = 0
	f.render()
}

func (f *feedViewer) render() {
	ui.Clear()
	ui.Render(f.grid)
}

func (f *feedViewer) initListView() {
	f.feedList = widgets.NewList()
	f.feedList.TextStyle = ui.NewStyle(ui.ColorWhite)
	f.feedList.BorderStyle.Fg = ui.ColorMagenta
	f.feedList.SelectedRowStyle.Fg = ui.ColorBlack
	f.feedList.SelectedRowStyle.Bg = ui.ColorWhite
}

func (f *feedViewer) initGrid() {
	f.grid = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	f.grid.SetRect(0, 0, termWidth, termHeight-1)

	f.infoText = widgets.NewParagraph()
	f.infoText.WrapText = true
	f.infoText.Border = false
	f.infoText.TextStyle = ui.Style{Modifier: ui.ModifierBold, Bg: ui.ColorWhite}
	f.infoText.Text = "initiating..."

	f.grid.Set(ui.NewRow(0.9, ui.NewCol(1.0, f.feedList)), ui.NewRow(0.1, ui.NewCol(1.0, f.infoText)))
}

func (f *feedViewer) openArticle() {
	index := f.feedList.SelectedRow
	if len(f.feed) > index {
		browser.Stderr = nil
		browser.Stdout = nil
		_ = browser.OpenURL(f.feed[index].Link)
		f.render()
	}
}

func (f *feedViewer) initEventsPolling() {
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			{
				switch e.ID {
				case "q", "C-c":
					return
				case "j":
					f.feedList.ScrollDown()
				case "k":
					f.feedList.ScrollUp()
				case "n":
					f.loadFeed(f.scraper.GetNextFeed)
				case "p":
					f.loadFeed(f.scraper.GetPrevFeed)
				case "r":
					f.loadFeed(f.scraper.GetInitialFeed)
				case "<Resize>":
					payload := e.Payload.(ui.Resize)
					f.grid.SetRect(0, 0, payload.Width, payload.Height-1)
				case "<Enter>":
					f.openArticle()
				}

				f.render()
			}
		}
	}
}
