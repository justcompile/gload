package internal

import (
	"net/http"
	"time"
)

type Options struct {
	Duration time.Duration
	Output   string
	Workers  int
}

type Result struct {
	Duration time.Duration
	Response *http.Response
}
