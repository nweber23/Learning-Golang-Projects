package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type CollectionStats struct {
	TotalIDsTried     int `json:"total_ids_tried"`
	SuccessfulFetches int `json:"successful_fetches"`
	UniquePDFsSaved   int `json:"unique_pdfs_saved"`
	DuplicatesFound   int `json:"duplicates_found"`
	Errors            int `json:"errors"`
	DurationSeconds   int `json:"duration_seconds"`
}

func savePDF(content []byte, id int, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("mkdir failed: %w", err)
	}
	filename := fmt.Sprintf("%d_en.subject.pdf", id)
	filePath := filepath.Join(outputDir, filename)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	return nil
}

func collectResults(
	resultChan chan FetchResult,
	config *ScraperConfig,
	done chan struct{},
) {
	deduper := NewDedupeTracker()
	uniqueSaved := 0
	duplicates := 0
	startTime := time.Now()
	for result := range resultChan {
		if deduper.HasSeen(result.ContentHash) {
			duplicates++
			continue
		}
		deduper.MarkSeen(result.ContentHash)
		if err := savePDF(result.Content, result.ID, config.OutputDir); err != nil {
			log.Printf("error saving PDF %d: %v", result.ID, err)
			continue
		}

		uniqueSaved++
	}
	stats := CollectionStats{
		TotalIDsTried:     config.IDEnd - config.IDStart + 1,
		SuccessfulFetches: uniqueSaved + duplicates,
		UniquePDFsSaved:   uniqueSaved,
		DuplicatesFound:   duplicates,
		DurationSeconds:   int(time.Since(startTime).Seconds()),
	}
	reportPath := filepath.Join(config.OutputDir, "deduplication_report.json")
	reportBytes, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		log.Printf("error marshaling report: %v", err)
	} else {
		if err := os.WriteFile(reportPath, reportBytes, 0644); err != nil {
			log.Printf("error writing report: %v", err)
		} else {
			log.Printf("Report saved to %s", reportPath)
		}
	}
	log.Printf("Collection complete: %d unique PDFs saved, %d duplicates removed", uniqueSaved, duplicates)
	done <- struct{}{}
}
