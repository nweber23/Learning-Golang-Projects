package main

import "time"

type FetchResult struct {
	ID          int
	Content     []byte
	ContentHash string
	FetchedAt   time.Time
}

type Stats struct {
	TotalTried      int
	SuccessfulFetch int
	DuplicatesFound int
	Errors          int
	Duration        time.Duration
}

type DedupeTracker struct {
	Hashes map[string]bool
}

func NewDedupeTracker() *DedupeTracker {
	return &DedupeTracker{
		Hashes: make(map[string]bool),
	}
}

func (dt *DedupeTracker) HasSeen(hash string) bool {
	return dt.Hashes[hash]
}

func (dt *DedupeTracker) MarkSeen(hash string) {
	dt.Hashes[hash] = true
}
