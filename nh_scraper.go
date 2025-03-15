package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/mmcdole/gofeed"
)

type nhscraper struct {
	output chan any
	log    *slog.Logger
}

func (h *nhscraper) parseFeed(ctx context.Context) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://news.ycombinator.com/rss")
	if err != nil {
		return err
	}
	for _, item := range feed.Items {
		post := hnpost{
			url:   item.Link,
			title: item.Title,
		}
		h.log.Debug("found item", slog.Any("item", *item))
		h.output <- &post
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
