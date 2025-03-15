package main

import "context"

type process interface {
	start(ctx context.Context) error
}
