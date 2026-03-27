package main

import (
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
