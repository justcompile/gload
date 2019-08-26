package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type Record struct {
	url      string
	status   int
	duration time.Duration
	stats    *Stats
}

type Stats struct {
	workers  int
	duration time.Duration
}

type Records map[string][]*Record

func ParseResults(inputFilePath string) (Records, error) {
	var reader io.Reader

	if inputFilePath == "-" {
		reader = os.Stdin
	} else {
		csvfile, err := os.Open(inputFilePath)
		if err != nil {
			return nil, err
		}

		reader = csvfile
	}

	// Parse the file
	r := csv.NewReader(reader)
	//r := csv.NewReader(bufio.NewReader(csvfile))

	results := make(map[string][]*Record)

	var stats *Stats

	// Iterate through the records
	for {
		// Read each record from csv
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if line[0] == "###STATS###" {
			stats = new(Stats)
			w := strings.Split(line[1], ":")

			noWorkers, _ := strconv.Atoi(w[1])
			stats.workers = noWorkers

			d := strings.Split(line[2], ":")
			dur, _ := time.ParseDuration(d[1])
			stats.duration = dur

			continue
		}

		status, err := strconv.Atoi(line[1])
		if err != nil {
			return nil, err
		}

		duration, err := time.ParseDuration(line[2])
		if err != nil {
			return nil, err
		}

		record := &Record{
			line[0],
			status,
			duration,
			stats,
		}

		if arr, hasKey := results[record.url]; hasKey {
			arr = append(arr, record)
			results[record.url] = arr
		} else {
			arr := make([]*Record, 1)
			arr[0] = record
			results[record.url] = arr
		}

	}

	return results, nil
}

func RenderResults(output io.Writer, records Records) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"URL", "Total Requests", "Requests p/s", "Min. Response Time", "Max. Response Time", "Avg. Response Time"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for url, results := range records {
		totalRequests := len(results)

		var totalResponseTime, minResponse, maxResponse int64
		minResponse = -1

		for _, result := range results {
			ns := result.duration.Nanoseconds()

			totalResponseTime += ns

			if ns < minResponse || minResponse == -1 {
				minResponse = ns
			}

			if ns > maxResponse {
				maxResponse = ns
			}
		}

		avgNs, err := time.ParseDuration(fmt.Sprintf("%dns", totalResponseTime/int64(totalRequests)))
		if err != nil {
			panic(err)
		}

		min, _ := time.ParseDuration(fmt.Sprintf("%dns", minResponse))
		max, _ := time.ParseDuration(fmt.Sprintf("%dns", maxResponse))

		reqPerSec := float64(totalRequests) / results[0].stats.duration.Seconds()

		table.Append([]string{
			url,
			fmt.Sprintf("%d", totalRequests),
			fmt.Sprintf("%f", reqPerSec),
			min.String(),
			max.String(),
			avgNs.String(),
		})
	}

	table.Render() // Send output
}
