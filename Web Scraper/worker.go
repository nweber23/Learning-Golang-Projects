package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func fetchPDF(url string, timeout time.Duration) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	return content, nil
}

func fetchWorker(
	idChan chan int,
	resultChan chan FetchResult,
	config *ScraperConfig,
) {
	for id := range idChan {
		url := fmt.Sprintf("%s/%d/en.subject.pdf", config.BaseURL, id)
		content, err := fetchPDF(url, config.Timeout)
		if err != nil {
			continue
		}
		hash := CalculateHash(content)
		result := FetchResult{
			ID:          id,
			Content:     content,
			ContentHash: hash,
			FetchedAt:   time.Now(),
		}
		resultChan <- result
	}
}
