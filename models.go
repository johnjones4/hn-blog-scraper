package main

import "time"

type hnpost struct {
	url   string
	title string
}

type scrapedSite struct {
	post            hnpost
	feedUrl         string
	siteTitle       string
	siteDescription string
	keywords        []string
	created         time.Time
}
