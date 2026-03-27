package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Test loading defaults from config.txt
	config, err := LoadConfig("config.txt")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify defaults
	if config.IDStart != 0 {
		t.Errorf("expected IDStart=0, got %d", config.IDStart)
	}
	if config.IDEnd != 1000000 {
		t.Errorf("expected IDEnd=1000000, got %d", config.IDEnd)
	}
	if config.Workers != 10 {
		t.Errorf("expected Workers=10, got %d", config.Workers)
	}
	if config.Timeout != 10*time.Second {
		t.Errorf("expected Timeout=10s, got %v", config.Timeout)
	}
}

func TestCalculateHash(t *testing.T) {
	content := []byte("test PDF content")
	hash := CalculateHash(content)

	// SHA256 produces exactly 64 hex characters
	if len(hash) != 64 {
		t.Errorf("expected 64 hex chars, got %d", len(hash))
	}

	// Hash should be lowercase hex
	for _, c := range hash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("expected hex chars, got %c", c)
		}
	}
}

func TestCalculateHash_Deterministic(t *testing.T) {
	content := []byte("test PDF content")
	hash1 := CalculateHash(content)
	hash2 := CalculateHash(content)

	// Same content should always produce the same hash
	if hash1 != hash2 {
		t.Errorf("hash should be deterministic: %s != %s", hash1, hash2)
	}
}

func TestCalculateHash_DifferentContent(t *testing.T) {
	content1 := []byte("PDF content A")
	content2 := []byte("PDF content B")

	hash1 := CalculateHash(content1)
	hash2 := CalculateHash(content2)

	// Different content should produce different hashes
	if hash1 == hash2 {
		t.Errorf("different content should produce different hash")
	}
}

func TestProduceIDs(t *testing.T) {
	idChan := make(chan int, 10) // Buffer of 10

	// Run producer in a goroutine
	go produceIDs(idChan, 0, 9)

	// Collect all IDs from the channel
	var ids []int
	for id := range idChan { // range reads until channel closes
		ids = append(ids, id)
	}

	// Should have 10 IDs (0-9 inclusive)
	if len(ids) != 10 {
		t.Errorf("expected 10 IDs, got %d", len(ids))
	}

	// Check they're in order
	for i := 0; i < 10; i++ {
		if ids[i] != i {
			t.Errorf("expected IDs in order, got %v", ids)
			break
		}
	}
}

func TestFetchPDF_Success(t *testing.T) {
	// Create a fake HTTP server that serves a PDF
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.Write([]byte("fake PDF content"))
	}))
	defer server.Close()

	// Fetch from the fake server
	content, err := fetchPDF(server.URL, 5*time.Second)
	if err != nil {
		t.Fatalf("fetchPDF failed: %v", err)
	}

	if string(content) != "fake PDF content" {
		t.Errorf("expected 'fake PDF content', got %q", string(content))
	}
}

func TestFetchPDF_404(t *testing.T) {
	// Create a server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := fetchPDF(server.URL, 5*time.Second)
	if err == nil {
		t.Errorf("expected error for 404, got nil")
	}
}

func TestFetchPDF_Timeout(t *testing.T) {
	// Create a server that sleeps longer than timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Sleep 2 seconds
	}))
	defer server.Close()

	// Timeout is 100ms, so this should fail
	_, err := fetchPDF(server.URL, 100*time.Millisecond)
	if err == nil {
		t.Errorf("expected timeout error, got nil")
	}
}

func TestFetchWorker(t *testing.T) {
	// Create a fake server that always returns the same PDF
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.Write([]byte("fake PDF"))
	}))
	defer server.Close()

	// Create channels
	idChan := make(chan int, 5)
	resultChan := make(chan FetchResult, 5)

	// Create config pointing to fake server
	config := &ScraperConfig{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}

	// WaitGroup to track worker completion
	var wg sync.WaitGroup

	// Start a worker goroutine
	wg.Add(1)
	go fetchWorker(idChan, resultChan, config, &wg)

	// Send 3 IDs
	idChan <- 1
	idChan <- 2
	idChan <- 3
	close(idChan) // Signal worker: no more IDs

	// Wait for worker to finish
	wg.Wait()

	// Close result channel so we can iterate
	close(resultChan)

	// Collect results
	var results []FetchResult
	for result := range resultChan {
		results = append(results, result)
	}

	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}

	// Check that hashes are calculated
	for _, result := range results {
		if result.ContentHash == "" {
			t.Errorf("expected non-empty hash, got empty")
		}
		if len(result.ContentHash) != 64 {
			t.Errorf("expected 64-char hash, got %d", len(result.ContentHash))
		}
	}
}

func TestSavePDF(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	content := []byte("fake PDF content")
	err := savePDF(content, 12345, tempDir)
	if err != nil {
		t.Fatalf("savePDF failed: %v", err)
	}

	// Check file was created
	filePath := filepath.Join(tempDir, "12345_en.subject.pdf")
	savedContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	if string(savedContent) != "fake PDF content" {
		t.Errorf("saved content mismatch: expected 'fake PDF content', got %q", string(savedContent))
	}
}

func TestDedupeTracker(t *testing.T) {
	deduper := NewDedupeTracker()

	hash1 := "abc123"
	hash2 := "xyz789"

	// Initially, no hashes are seen
	if deduper.HasSeen(hash1) {
		t.Errorf("expected HasSeen(hash1) = false initially")
	}

	// Mark hash1 as seen
	deduper.MarkSeen(hash1)
	if !deduper.HasSeen(hash1) {
		t.Errorf("expected HasSeen(hash1) = true after MarkSeen")
	}

	// hash2 should still not be seen
	if deduper.HasSeen(hash2) {
		t.Errorf("expected HasSeen(hash2) = false")
	}

	// Mark hash2 as seen
	deduper.MarkSeen(hash2)
	if !deduper.HasSeen(hash2) {
		t.Errorf("expected HasSeen(hash2) = true after MarkSeen")
	}
}

func TestCollectResults(t *testing.T) {
	tempDir := t.TempDir()

	resultChan := make(chan FetchResult, 5)
	done := make(chan struct{})

	config := &ScraperConfig{
		IDStart:   0,
		IDEnd:     5,
		OutputDir: tempDir,
	}

	// Start collector goroutine
	go collectResults(resultChan, config, done)

	// Send 3 results with same content (2 duplicates)
	sameContent := []byte("identical PDF")
	sameHash := CalculateHash(sameContent)

	resultChan <- FetchResult{ID: 1, Content: sameContent, ContentHash: sameHash, FetchedAt: time.Now()}
	resultChan <- FetchResult{ID: 2, Content: sameContent, ContentHash: sameHash, FetchedAt: time.Now()}
	resultChan <- FetchResult{ID: 3, Content: []byte("different PDF"), ContentHash: CalculateHash([]byte("different PDF")), FetchedAt: time.Now()}

	close(resultChan) // Signal no more results

	// Wait for collector to finish
	<-done

	// Check files were saved (should be 2 unique)
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("failed to read output dir: %v", err)
	}

	pdfCount := 0
	for _, f := range files {
		if f.Name() != "deduplication_report.json" {
			pdfCount++
		}
	}

	if pdfCount != 2 {
		t.Errorf("expected 2 unique PDFs (1 duplicate removed), got %d", pdfCount)
	}

	// Check report was created
	reportPath := filepath.Join(tempDir, "deduplication_report.json")
	if _, err := os.Stat(reportPath); err != nil {
		t.Errorf("expected deduplication_report.json to exist")
	}
}

func TestFullPipeline(t *testing.T) {
	// Create a mock server that returns PDFs
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.Write([]byte("fake PDF"))
	}))
	defer server.Close()

	// Create temp directory for output
	tempDir := t.TempDir()

	// Create config pointing to mock server
	config := &ScraperConfig{
		BaseURL:          server.URL,
		IDStart:          0,
		IDEnd:            5,
		OutputDir:        tempDir,
		Workers:          2,
		Timeout:          5 * time.Second,
		BufferIDChan:     10,
		BufferResultChan: 5,
	}

	// Create channels
	idChan := make(chan int, config.BufferIDChan)
	resultChan := make(chan FetchResult, config.BufferResultChan)
	done := make(chan struct{})

	var wg sync.WaitGroup

	// Start producer
	go produceIDs(idChan, config.IDStart, config.IDEnd)

	// Start workers
	for i := 0; i < config.Workers; i++ {
		wg.Add(1)
		go fetchWorker(idChan, resultChan, config, &wg)
	}

	// Start collector
	go collectResults(resultChan, config, done)

	// Wait for all workers to finish
	wg.Wait()

	// Close result channel
	close(resultChan)

	// Wait for collector to finish
	<-done

	// Check results
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("failed to read output dir: %v", err)
	}

	// Should have 1 unique PDF (all are identical) + report
	pdfCount := 0
	for _, f := range files {
		if f.Name() != "deduplication_report.json" {
			pdfCount++
		}
	}

	if pdfCount != 1 {
		t.Errorf("expected 1 unique PDF (5 duplicates removed), got %d", pdfCount)
	}

	// Verify report exists
	reportPath := filepath.Join(tempDir, "deduplication_report.json")
	if _, err := os.Stat(reportPath); err != nil {
		t.Errorf("expected deduplication_report.json, got error: %v", err)
	}
}
