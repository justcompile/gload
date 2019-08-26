package internal

import (
	"net/http"
	"time"
)

type Options struct {
	Duration time.Duration
	Workers  int
}

type Result struct {
	Duration time.Duration
	Response *http.Response
}
