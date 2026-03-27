package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// fetchPDF fetches a PDF from the given URL with timeout.
// Returns raw bytes or error.
func fetchPDF(url string, timeout time.Duration) ([]byte, error) {
	// Create a client with timeout
	client := &http.Client{
		Timeout: timeout,
	}

	// Make GET request
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Read response body into memory
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	return content, nil
}

// fetchWorker reads IDs from idChan, fetches PDFs, sends results to resultChan.
// Calls wg.Done() when it exits to signal completion to main.
func fetchWorker(
	idChan chan int,
	resultChan chan FetchResult,
	config *ScraperConfig,
	wg *sync.WaitGroup,
) {
	// Signal completion when this goroutine exits
	// defer ensures this runs even if there's an error
	defer wg.Done()

	processed := 0

	// Keep reading IDs until channel closes
	for id := range idChan {
		processed++

		// Print progress every 1000 IDs attempted
		if processed%1000 == 0 {
			log.Printf("[WORKER] Attempted %d IDs so far (current ID: %d)", processed, id)
		}

		// Build URL for this ID
		url := fmt.Sprintf("%s/%d/en.subject.pdf", config.BaseURL, id)

		// Fetch the PDF
		content, err := fetchPDF(url, config.Timeout)
		if err != nil {
			// Log error silently and continue to next ID
			continue
		}

		// Calculate hash of the PDF content
		hash := CalculateHash(content)

		// Create result and send down result channel
		result := FetchResult{
			ID:          id,
			Content:     content,
			ContentHash: hash,
			FetchedAt:   time.Now(),
		}
		resultChan <- result
	}

	// When idChan closes, loop exits, defer runs wg.Done()
}
