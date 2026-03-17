package models

import "time"

type Transformation struct {
	ID              string                 `json:"id"`
	UserID          string                 `json:"user_id"`
	OriginalImageID string                 `json:"original_image_id"`
	Transformations map[string]interface{} `json:"transformations"`
	ResultURL       string                 `json:"result_url"`
	CreatedAt       time.Time              `json:"created_at"`
}

type TransformationRequest struct {
	Transformations TransformOptions `json:"transformations"`
}

type TransformOptions struct {
	Resize  *ResizeOptions `json:"resize,omitempty"`
	Crop    *CropOptions   `json:"crop,omitempty"`
	Rotate  int            `json:"rotate,omitempty"`
	Flip    string         `json:"flip,omitempty"`
	Filters *FilterOptions `json:"filters,omitempty"`
	Format  string         `json:"format,omitempty"`
}

type ResizeOptions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type CropOptions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type FilterOptions struct {
	Grayscale  bool    `json:"grayscale,omitempty"`
	Sepia      bool    `json:"sepia,omitempty"`
	Blur       float64 `json:"blur,omitempty"`
	Brightness float64 `json:"brightness,omitempty"`
	Contrast   float64 `json:"contrast,omitempty"`
}
