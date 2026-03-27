package main

import (
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
