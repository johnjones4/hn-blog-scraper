package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
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
	db, err := sql.Open("sqlite3", os.Getenv("DB_PATH"))
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
	defer res.Close()

	return res.Next(), nil
}

func (s *store) insertScrape(ctx context.Context, site *scrapedSite) error {
	keywords, err := json.Marshal(site.Keywords)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, "insert into scraped_site (post_url, post_title, feed_url, site_title, site_description, keywords, created) values (?,?,?,?,?,?,?)",
		site.Post.Url,
		site.Post.Title,
		site.FeedUrl,
		site.SiteTitle,
		site.SiteDescription,
		string(keywords),
		site.Created.Format(time.RFC3339Nano),
	)
	if err != nil {
		return err
	}
	s.log.Info("inserted site", slog.String("title", site.SiteTitle))
	return nil
}

func (s *store) getScraped(ctx context.Context) ([]*scrapedStat, error) {
	var sites []*scrapedStat
	rows, err := s.db.QueryContext(ctx, "SELECT site_title, site_description, feed_url, COUNT(DISTINCT post_url) AS unique_post_count FROM scraped_site GROUP BY feed_url order by COUNT(DISTINCT post_url) desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var stat scrapedStat
		err = rows.Scan(
			&stat.SiteTitle,
			&stat.SiteDescription,
			&stat.FeedUrl,
			&stat.PostCount,
		)
		if err != nil {
			return nil, err
		}
		sites = append(sites, &stat)
	}
	return sites, nil
}
