package internal

import (
	"context"
	"fmt"
)

func worker(ctx context.Context, url string, results chan<- *Result) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			results <- makeRequest(url)
		}
	}
}

type Dispatcher struct {
	opts *Options
}

func (d *Dispatcher) Run(url string) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, d.opts.Duration)
	defer cancel()

	results := make(chan *Result, 1000)

	for i := 0; i <= d.opts.Workers; i++ {
		go worker(ctx, url, results)
	}

	defer close(results)

	for {
		select {
		case r := <-results:
			fmt.Println(r)
		case <-ctx.Done():
			return
		}
	}
}

func NewDispatcher(opts *Options) *Dispatcher {
	return &Dispatcher{opts}
}
