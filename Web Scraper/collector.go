package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// CollectionStats represents the final deduplication report
type CollectionStats struct {
	TotalIDsTried     int `json:"total_ids_tried"`
	SuccessfulFetches int `json:"successful_fetches"`
	UniquePDFsSaved   int `json:"unique_pdfs_saved"`
	DuplicatesFound   int `json:"duplicates_found"`
	Errors            int `json:"errors"`
	DurationSeconds   int `json:"duration_seconds"`
}

// savePDF writes PDF content to disk with a safe filename
func savePDF(content []byte, id int, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("mkdir failed: %w", err)
	}

	// Use ID as filename (original filename is always en.subject.pdf)
	filename := fmt.Sprintf("%d_en.subject.pdf", id)
	filePath := filepath.Join(outputDir, filename)

	// Write file to disk
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	return nil
}

// collectResults reads from resultChan, deduplicates by content hash, and saves PDFs.
// Runs in its own goroutine. Signals done when finished.
func collectResults(
	resultChan chan FetchResult,
	config *ScraperConfig,
	done chan struct{},
) {
	deduper := NewDedupeTracker()
	uniqueSaved := 0
	duplicates := 0

	startTime := time.Now()

	// Read results until channel closes
	for result := range resultChan {
		// Check if we've seen this hash before
		if deduper.HasSeen(result.ContentHash) {
			// This is a duplicate - skip it
			duplicates++
			continue
		}

		// Mark this hash as seen
		deduper.MarkSeen(result.ContentHash)

		// Try to save the PDF to disk
		if err := savePDF(result.Content, result.ID, config.OutputDir); err != nil {
			log.Printf("error saving PDF %d: %v", result.ID, err)
			continue
		}

		uniqueSaved++
	}

	// All results received, finalize statistics
	stats := CollectionStats{
		TotalIDsTried:     config.IDEnd - config.IDStart + 1,
		SuccessfulFetches: uniqueSaved + duplicates,
		UniquePDFsSaved:   uniqueSaved,
		DuplicatesFound:   duplicates,
		DurationSeconds:   int(time.Since(startTime).Seconds()),
	}

	// Save deduplication report as JSON
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

	// Signal to main that we're done
	done <- struct{}{}
}
