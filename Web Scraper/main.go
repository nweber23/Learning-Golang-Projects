package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	// Load configuration from config.txt + CLI flags
	config, err := LoadConfig("config.txt")
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	log.Printf("Starting scraper: %d workers, IDs %d to %d",
		config.Workers, config.IDStart, config.IDEnd)

	// Create channels
	idChan := make(chan int, config.BufferIDChan)
	resultChan := make(chan FetchResult, config.BufferResultChan)
	done := make(chan struct{})

	startTime := time.Now()

	// WaitGroup tracks when all workers finish
	var wg sync.WaitGroup

	// Start producer goroutine
	go produceIDs(idChan, config.IDStart, config.IDEnd)

	// Start N worker goroutines
	for i := 0; i < config.Workers; i++ {
		wg.Add(1)  // Say "one more worker starting"
		go fetchWorker(idChan, resultChan, config, &wg)
	}

	// Start collector goroutine
	go collectResults(resultChan, config, done)

	// Wait for all workers to finish
	// Each worker calls wg.Done() when it exits
	// This blocks until all N workers have called Done()
	wg.Wait()

	// Now we know all workers are done
	// Safe to close resultChan so collector knows no more results coming
	close(resultChan)

	// Wait for collector to finish processing and save report
	<-done

	duration := time.Since(startTime)
	log.Printf("Scraping completed in %v", duration)
}
