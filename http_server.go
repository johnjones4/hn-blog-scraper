package main

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "embed"
)

//go:embed index.html
var index string

type httpServer struct {
	log   *slog.Logger
	store *store
}

var tpl *template.Template

func (s *httpServer) start(ctx context.Context) error {
	return http.ListenAndServe(os.Getenv("HTTP_HOST"), http.HandlerFunc(s.handle))
}

func (s *httpServer) handle(w http.ResponseWriter, r *http.Request) {
	if tpl == nil {
		var err error
		tpl, err = template.New("index").Parse(index)
		if err != nil {
			s.log.Error("error loading template", slog.Any("error", err), slog.String("template", index))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	sites, err := s.store.getScraped(r.Context())
	if err != nil {
		s.log.Error("error getting sites", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, struct {
		Sites []*scrapedStat
	}{
		Sites: sites,
	})
	if err != nil {
		s.log.Error("error executing template", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
