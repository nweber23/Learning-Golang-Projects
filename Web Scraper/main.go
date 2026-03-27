package main

import (
	"log"
	"time"
)

func main() {
	config, err := LoadConfig("config.txt")
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	log.Printf("Starting scraper: %d workers, IDs %d to %d",
		config.Workers, config.IDStart, config.IDEnd)
	idChan := make(chan int, config.BufferIDChan)
	resultChan := make(chan FetchResult, config.BufferResultChan)
	done := make(chan struct{})
	startTime := time.Now()
	go produceIDs(idChan, config.IDStart, config.IDEnd)
	for i := 0; i < config.Workers; i++ {
		go fetchWorker(idChan, resultChan, config)
	}
	go collectResults(resultChan, config, done)
	<-done
	duration := time.Since(startTime)
	log.Printf("Scraping completed in %v", duration)
}
