package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("config.txt")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
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
	if len(hash) != 64 {
		t.Errorf("expected 64 hex chars, got %d", len(hash))
	}
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
	if hash1 != hash2 {
		t.Errorf("hash should be deterministic: %s != %s", hash1, hash2)
	}
}

func TestCalculateHash_DifferentContent(t *testing.T) {
	content1 := []byte("PDF content A")
	content2 := []byte("PDF content B")
	hash1 := CalculateHash(content1)
	hash2 := CalculateHash(content2)
	if hash1 == hash2 {
		t.Errorf("different content should produce different hash")
	}
}

func TestProduceIDs(t *testing.T) {
	idChan := make(chan int, 10)
	go produceIDs(idChan, 0, 9)
	var ids []int
	for id := range idChan {
		ids = append(ids, id)
	}
	if len(ids) != 10 {
		t.Errorf("expected 10 IDs, got %d", len(ids))
	}
	for i := 0; i < 10; i++ {
		if ids[i] != i {
			t.Errorf("expected IDs in order, got %v", ids)
			break
		}
	}
}

func TestFetchPDF_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.Write([]byte("fake PDF content"))
	}))
	defer server.Close()
	content, err := fetchPDF(server.URL, 5*time.Second)
	if err != nil {
		t.Fatalf("fetchPDF failed: %v", err)
	}
	if string(content) != "fake PDF content" {
		t.Errorf("expected 'fake PDF content', got %q", string(content))
	}
}

func TestFetchPDF_404(t *testing.T) {
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()
	_, err := fetchPDF(server.URL, 100*time.Millisecond)
	if err == nil {
		t.Errorf("expected timeout error, got nil")
	}
}

func TestFetchWorker(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.Write([]byte("fake PDF"))
	}))
	defer server.Close()
	idChan := make(chan int, 5)
	resultChan := make(chan FetchResult, 5)
	config := &ScraperConfig{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	go fetchWorker(idChan, resultChan, config)
	idChan <- 1
	idChan <- 2
	idChan <- 3
	close(idChan)
	var results []FetchResult
	for result := range resultChan {
		results = append(results, result)
		if len(results) == 3 {
			break
		}
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
	for _, result := range results {
		if result.ContentHash == "" {
			t.Errorf("expected non-empty hash, got empty")
		}
		if len(result.ContentHash) != 64 {
			t.Errorf("expected 64-char hash, got %d", len(result.ContentHash))
		}
	}
}
