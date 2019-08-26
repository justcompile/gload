package internal

import (
	"log"
	"net/http"
	"time"
)

func makeRequest(url string) *Result {
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	client := http.Client{
		Timeout: time.Second * 5,
	}

	start := time.Now()

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	return &Result{
		Duration: time.Now().Sub(start),
		Response: resp,
	}
}
