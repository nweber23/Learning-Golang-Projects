package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ScraperConfig holds all configuration for the scraper
type ScraperConfig struct {
	BaseURL          string        // Base URL pattern (e.g., https://cdn.intra.42.fr/pdf/pdf)
	IDStart          int           // Start of ID range
	IDEnd            int           // End of ID range
	OutputDir        string        // Directory to save PDFs
	Workers          int           // Number of worker goroutines to run concurrently
	Timeout          time.Duration // Per-request HTTP timeout
	BufferIDChan     int           // Buffer size for ID channel (backpressure control)
	BufferResultChan int           // Buffer size for result channel
}

// LoadConfig reads config.txt and applies CLI flag overrides
func LoadConfig(filename string) (*ScraperConfig, error) {
	// Start with defaults
	config := &ScraperConfig{
		BaseURL:          "https://cdn.intra.42.fr/pdf/pdf",
		IDStart:          0,
		IDEnd:            1000000,
		OutputDir:        "./scraped_pdfs",
		Workers:          10,
		Timeout:          10 * time.Second,
		BufferIDChan:     100,
		BufferResultChan: 50,
	}

	// Try to load from config file
	if _, err := os.Stat(filename); err == nil {
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// Skip empty lines and comments
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			// Parse KEY=VALUE
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Map config values to struct fields
			switch key {
			case "BASE_URL":
				config.BaseURL = value
			case "ID_START":
				if v, err := strconv.Atoi(value); err == nil {
					config.IDStart = v
				}
			case "ID_END":
				if v, err := strconv.Atoi(value); err == nil {
					config.IDEnd = v
				}
			case "OUTPUT_DIR":
				config.OutputDir = value
			case "WORKERS":
				if v, err := strconv.Atoi(value); err == nil {
					config.Workers = v
				}
			case "TIMEOUT":
				if v, err := time.ParseDuration(value); err == nil {
					config.Timeout = v
				}
			case "BUFFER_SIZE_IDS":
				if v, err := strconv.Atoi(value); err == nil {
					config.BufferIDChan = v
				}
			case "BUFFER_SIZE_RESULTS":
				if v, err := strconv.Atoi(value); err == nil {
					config.BufferResultChan = v
				}
			}
		}
	}

	// Apply CLI flag overrides (these take precedence over config file)
	flag.IntVar(&config.IDStart, "id-start", config.IDStart, "Start ID")
	flag.IntVar(&config.IDEnd, "id-end", config.IDEnd, "End ID")
	flag.StringVar(&config.OutputDir, "output-dir", config.OutputDir, "Output directory")
	flag.IntVar(&config.Workers, "workers", config.Workers, "Number of workers")
	flag.DurationVar(&config.Timeout, "timeout", config.Timeout, "Request timeout")
	flag.Parse()

	// Validate configuration
	if config.IDStart < 0 || config.IDEnd <= config.IDStart {
		return nil, fmt.Errorf("invalid ID range: %d to %d", config.IDStart, config.IDEnd)
	}
	if config.Workers <= 0 {
		return nil, fmt.Errorf("workers must be > 0, got %d", config.Workers)
	}

	return config, nil
}
