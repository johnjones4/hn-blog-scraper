package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/mmcdole/gofeed"
)

var feeds = []string{
	"https://twostopbits.com/rss",
	"https://news.ycombinator.com/rss",
}

type nhscraper struct {
	output chan any
	log    *slog.Logger
}

func (h *nhscraper) parseFeed(ctx context.Context) error {
	for _, url := range feeds {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(url)
		if err != nil {
			h.log.Error("error parsing feed", slog.Any("error", err), slog.String("url", url))
		}
		for _, item := range feed.Items {
			post := hnpost{
				url:   item.Link,
				title: item.Title,
			}
			h.log.Info("found item", slog.Any("item", *item))
			h.output <- &post
		}
	}
	return nil
}

func (h *nhscraper) start(ctx context.Context) error {
	defer close(h.output)
	err := h.parseFeed(ctx)
	if err != nil {
		return err
	}
	tick := time.Tick(time.Minute * 5)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick:
			err := h.parseFeed(ctx)
			if err != nil {
				h.log.Error("error parsing hn feed", slog.Any("error", err))
			}
		}
	}
}
