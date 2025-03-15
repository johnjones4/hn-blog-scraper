package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schema string

type store struct {
	db  *sql.DB
	in  chan any
	log *slog.Logger
}

func (s *store) start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case scrapedAny := <-s.in:
			scraped := scrapedAny.(*scrapedSite)
			err := s.insertScrape(ctx, scraped)
			if err != nil {
				s.log.Error("error inserting scraped site", slog.Any("error", err))
			}
		}
	}
}

func (s *store) init(ctx context.Context) error {
	db, err := sql.Open("sqlite3", "./data.sqlite")
	if err != nil {
		return err
	}
	s.db = db
	_, err = s.db.ExecContext(ctx, schema)
	return err
}

func (s *store) hasPostBeenScraped(ctx context.Context, postUrl string) (bool, error) {
	res, err := s.db.QueryContext(ctx, "select * from scraped_site where post_url = ?", postUrl)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	return res.Next(), nil
}

func (s *store) insertScrape(ctx context.Context, site *scrapedSite) error {
	_, err := s.db.ExecContext(ctx, "insert into scraped_site (post_url, post_title, feed_url, site_title, site_description, created) values (?,?,?,?,?,?)",
		site.post.url,
		site.post.title,
		site.feedUrl,
		site.siteTitle,
		site.siteDescription,
		time.Now(),
	)
	return err
}
