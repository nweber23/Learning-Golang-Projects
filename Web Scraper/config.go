package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type ScraperConfig struct {
	BaseURL          string
	IDStart          int
	IDEnd            int
	OutputDir        string
	Workers          int
	Timeout          time.Duration
	BufferIDChan     int
	BufferResultChan int
}

func LoadConfig(filename string) (*ScraperConfig, error) {
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
	if _, err := os.Stat(filename); err == nil {
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %s", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Printf("failed to close config file: %s", err)
			}
		}()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			switch key {
			case "BaseURL":
				config.BaseURL = value
			case "IDStart":
				if value, err := strconv.Atoi(value); err == nil {
					config.IDStart = value
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
	flag.IntVar(&config.IDStart, "id-start", config.IDStart, "Start ID")
	flag.IntVar(&config.IDEnd, "id-end", config.IDEnd, "End ID")
	flag.StringVar(&config.OutputDir, "output-dir", config.OutputDir, "Output directory")
	flag.IntVar(&config.Workers, "workers", config.Workers, "Number of workers")
	flag.DurationVar(&config.Timeout, "timeout", config.Timeout, "Request timeout")
	flag.Parse()
	if config.IDStart < 0 || config.IDEnd <= config.IDStart {
		return nil, fmt.Errorf("invalid ID range: %d to %d", config.IDStart, config.IDEnd)
	}
	if config.Workers <= 0 {
		return nil, fmt.Errorf("workers must be > 0, got %d", config.Workers)
	}
	return config, nil
}
