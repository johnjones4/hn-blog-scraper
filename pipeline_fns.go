package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	rake "github.com/afjoseph/RAKE.go"
	"github.com/microcosm-cc/bluemonday"
)

func scrapeSite(post *hnpost) (*scrapedSite, error) {
	res, err := http.Get(post.url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", res.StatusCode)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	sel := doc.Find("link[type=\"application/rss+xml\"]").First()

	if sel == nil {
		return nil, nil
	}

	href := sel.AttrOr("href", "")
	if href == "" {
		return nil, nil
	}

	parsedUrl, err := url.Parse(post.url)
	if err != nil {
		return nil, err
	}

	parsedFeedUrl, err := url.Parse(href)
	if err != nil {
		return nil, err
	}

	resolvedParsedUrl := parsedUrl.ResolveReference(parsedFeedUrl)

	titleEl := doc.Find("title").First()
	var title string
	if titleEl != nil {
		title = strings.TrimSpace(titleEl.Text())
	}

	descEl := doc.Find("meta[name=\"description\"]").First()
	var description string
	if descEl != nil {
		description = strings.TrimSpace(descEl.AttrOr("content", ""))
	}

	body := doc.Find("body").First().Text()
	cleanBody := bluemonday.StrictPolicy().Sanitize(body)
	candidates := rake.RunRake(cleanBody)
	keywords := make([]string, 0, len(candidates))
	for _, c := range candidates {
		keywords = append(keywords, c.Key)
	}

	return &scrapedSite{
		post:            *post,
		feedUrl:         resolvedParsedUrl.String(),
		siteTitle:       title,
		siteDescription: description,
		created:         time.Now(),
		keywords:        keywords,
	}, nil
}
