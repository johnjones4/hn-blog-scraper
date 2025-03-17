package main

import "time"

type hnpost struct {
	Url   string
	Title string
}

type scrapedSite struct {
	Post            hnpost
	FeedUrl         string
	SiteTitle       string
	SiteDescription string
	Keywords        []string
	Created         time.Time
}

type scrapedStat struct {
	SiteTitle       string
	SiteDescription string
	FeedUrl         string
	PostCount       int
}
