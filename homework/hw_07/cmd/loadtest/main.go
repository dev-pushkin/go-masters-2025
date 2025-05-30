package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go_course_master/homework/hw_07/pkg/loadtest"
)

// main is the entry point for the load testing application.
// Example usage:
// go run main.go http://localhost:8082 100000 100 10

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run main.go <url> <requests> <concurrency>")
		return
	}

	url := os.Args[1]
	totalRequests := mustConvertToInt(os.Args[2])
	concurrency := mustConvertToInt(os.Args[3])
	timeout := mustConvertToInt(os.Args[4])

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	result := loadtest.Run(ctx, url, totalRequests, concurrency)
	printReport(result)

}

func mustConvertToInt(s string) int {
	n, err := fmt.Sscanf(s, "%d", new(int))
	if err != nil || n != 1 {
		panic("invalid param, must be number: " + s)
	}
	var val int
	fmt.Sscanf(s, "%d", &val)
	return val
}

func printReport(result loadtest.TestResult) {
	fmt.Println("===== Load Test Report =====")
	fmt.Printf("Target URL: %s\n", result.TargetURL)
	fmt.Printf("Total requests: %d\n", result.TotalRequests)
	fmt.Printf("Concurrency: %d\n", result.Concurrency)
	fmt.Printf("Duration: %s\n", result.Duration)
	fmt.Printf("Success: %d\n", result.Success)
	fmt.Printf("Fail: %d\n", result.Fail)
	if result.Success > 0 {
		fmt.Printf("Avg Latency: %s\n", result.AvgLatency)
	}
	fmt.Printf("Max Latency: %s\n", result.MaxLatency)
	if result.Error != nil {
		fmt.Printf("\033[31mError: %s\033[0m\n", result.Error) // Красный
	} else {
		fmt.Printf("\033[32mNo errors encountered during the test.\033[0m\n") // Зелёный
	}
}
