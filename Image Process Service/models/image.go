package models

import "time"

type Image struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Filename    string    `json:"filename"`
	OriginalURL string    `json:"original_url"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Size        int64     `json:"size"`
	Format      string    `json:"format"`
	CreatedAt   time.Time `json:"created_at"`
}

type ImageResponse struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	URL       string    `json:"url"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Size      int64     `json:"size"`
	Format    string    `json:"format"`
	CreatedAt time.Time `json:"created_at"`
}

type ImageRequest struct {
	Images []*ImageResponse `json:"images"`
	Page   int              `json:"page"`
	Limit  int              `json:"limit"`
	Total  int              `json:"total"`
}
