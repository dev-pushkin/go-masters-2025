package loadtest

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	DefaultMaxIdleConns        = 1000000
	DefaultMaxIdleConnsPerHost = 1000000
	DefaultIdleConnTimeout     = 90 * time.Second
)

type reqResult struct {
	Latency time.Duration
	Error   error
}

type TestResult struct {
	TargetURL     string
	TotalRequests int
	Concurrency   int
	Duration      time.Duration
	Success       int
	Fail          int
	AvgLatency    time.Duration
	MaxLatency    time.Duration
	Error         error
}

func worker(client *http.Client, url string, jobs <-chan int, results chan<- reqResult) {
	for range jobs {
		start := time.Now()
		resp, err := client.Get(url)
		latency := time.Since(start)

		if err != nil {
			results <- reqResult{Latency: latency, Error: err}
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		results <- reqResult{Latency: latency, Error: nil}
	}
}

func Run(ctx context.Context, url string, totalRequests int, concurrency int) TestResult {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:      DefaultMaxIdleConns,
			MaxConnsPerHost:   DefaultMaxIdleConnsPerHost,
			IdleConnTimeout:   DefaultIdleConnTimeout,
			DisableKeepAlives: false,
		},
	}

	jobs := make(chan int)
	results := make(chan reqResult, totalRequests)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(client, url, jobs, results)
		}()
	}

	startTime := time.Now()

	go func() {
		defer close(jobs)
		for i := range totalRequests {
			select {
			case <-ctx.Done():
				return // Exit if context is done
			case jobs <- i:
				// Send job to channel
			}
		}
	}()

	wg.Wait()
	close(results)

	duration := time.Since(startTime)

	var totalLatency time.Duration
	var maxLatency time.Duration
	var success, fail int
	var err error

DONE:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			break DONE
		case r, ok := <-results:
			if !ok {
				// If results channel is closed, break the loop
				break DONE
			}
			if r.Error != nil {
				fail++
				continue
			}
			success++
			totalLatency += r.Latency
			if r.Latency > maxLatency {
				maxLatency = r.Latency
			}

		}
	}

	return TestResult{
		TargetURL:     url,
		TotalRequests: totalRequests,
		Concurrency:   concurrency,
		Duration:      duration,
		Success:       success,
		Fail:          fail,
		AvgLatency:    totalLatency / time.Duration(success),
		MaxLatency:    maxLatency,
		Error:         err,
	}
}
