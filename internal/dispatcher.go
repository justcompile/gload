package internal

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
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
	opts   *Options
	stream *bufio.Writer
	f      io.Writer
}

func (d *Dispatcher) Close() error {
	return d.stream.Flush()
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

	d.stream.WriteString("###STATS###,")
	d.stream.WriteString(fmt.Sprintf("workers:%d", d.opts.Workers))
	d.stream.WriteRune(',')
	d.stream.WriteString("duration:")
	d.stream.WriteString(d.opts.Duration.String())
	d.stream.WriteString("\n")

	for {
		select {
		case r := <-results:
			d.render(r)
		case <-ctx.Done():
			return
		}
	}
}

func (d *Dispatcher) render(res *Result) {
	d.stream.WriteString(res.Response.Request.URL.String())
	d.stream.WriteRune(',')
	d.stream.WriteString(strconv.Itoa(res.Response.StatusCode))
	d.stream.WriteRune(',')
	d.stream.WriteString(res.Duration.String())
	d.stream.WriteString("\n")
}

func NewDispatcher(opts *Options) *Dispatcher {
	var stream io.Writer
	if opts.Output == "-" {
		stream = os.Stdout
	} else {
		f, err := os.Create(opts.Output)
		if err != nil {
			panic(err)
		}

		stream = f
	}

	return &Dispatcher{
		opts,
		bufio.NewWriter(stream),
		stream,
	}
}
