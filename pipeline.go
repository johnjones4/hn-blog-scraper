package main

import (
	"context"
	"log/slog"

	"github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"
)

type pipeline struct {
	in    chan any
	out   chan any
	store *store
	log   *slog.Logger
}

func (p *pipeline) start(ctx context.Context) error {
	extension.NewChanSource(p.in).
		Via(flow.NewFilter[*hnpost](func(h *hnpost) bool {
			exists, err := p.store.hasPostBeenScraped(ctx, h.url)
			if err != nil {
				p.log.Error("error checking post status", slog.Any("error", err))
				return false
			}
			return !exists
		}, 10)).
		Via(flow.NewMap[*hnpost, *scrapedSite](func(h *hnpost) *scrapedSite {
			scraped, err := scrapeSite(h)
			if err != nil {
				p.log.Error("error scraping site", slog.Any("error", err))
				return nil
			}
			return scraped
		}, 10)).
		Via(flow.NewFilter[*scrapedSite](func(ss *scrapedSite) bool {
			return ss != nil
		}, 10)).
		To(extension.NewChanSink(p.out))
	return nil
}
