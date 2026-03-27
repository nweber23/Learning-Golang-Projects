package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// FetchResult represents a successful PDF fetch from a worker
// Workers will send this down the result channel
type FetchResult struct {
	ID           int       // The ID that was fetched (e.g., 193946)
	Content      []byte    // Raw PDF bytes
	ContentHash  string    // SHA256 hex digest (64 chars) - used for deduplication
	FetchedAt    time.Time // When it was fetched
}

// Stats tracks scraping progress and completion metrics
type Stats struct {
	TotalTried      int
	SuccessfulFetch int
	DuplicatesFound int
	Errors          int
	Duration        time.Duration
}

// DedupeTracker detects duplicate PDFs by content hash
// The collector goroutine uses this to track which PDFs we've already seen
type DedupeTracker struct {
	Hashes map[string]bool // hash -> true if seen
}

// NewDedupeTracker creates an empty tracker
func NewDedupeTracker() *DedupeTracker {
	return &DedupeTracker{
		Hashes: make(map[string]bool),
	}
}

// HasSeen returns true if this hash was already seen
func (dt *DedupeTracker) HasSeen(hash string) bool {
	return dt.Hashes[hash]
}

// MarkSeen adds a hash to the seen set
func (dt *DedupeTracker) MarkSeen(hash string) {
	dt.Hashes[hash] = true
}

// CalculateHash returns the SHA256 hex digest of content
// Used for deduplication - identical PDFs will have identical hashes
func CalculateHash(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("%x", hash)
}
