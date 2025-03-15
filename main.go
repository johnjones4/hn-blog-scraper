package main

import (
	"context"
	"log/slog"
	"sync"
)

func main() {
	log := slog.Default()

	nhScraper := &nhscraper{
		output: make(chan any),
		log:    log,
	}

	s := &store{
		in:  make(chan any),
		log: log,
	}

	err := s.init(context.Background())
	if err != nil {
		panic(err)
	}

	pipe := &pipeline{
		in:    nhScraper.output,
		out:   s.in,
		log:   log,
		store: s,
	}

	processes := []process{
		nhScraper,
		s,
		pipe,
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	for _, p := range processes {
		wg.Add(1)
		go func() {
			err := p.start(ctx)
			if err != nil {
				log.Error("error running process", slog.Any("error", err))
				cancel()
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
