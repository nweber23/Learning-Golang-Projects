package processor

import (
	"errors"
	"image"

	"image-process-service/models"
)

type Processor struct{}

func newProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Process(src image.Image, req *models.TransformOptions) (image.Image, error) {
	if src == nil {
		return nil, errors.New("src can't be nil")
	}
	result := src
	if req.Resize != nil {
		result = Resize(result, req.Resize.Width, req.Resize.Height)
	}
	if req.Crop != nil {
		result = Crop(result, req.Crop.X, req.Crop.Y, req.Crop.Width, req.Crop.Height)
	}
	if req.Rotate != 0 {
		result = Rotate(result, req.Rotate)
	}
	if req.Flip != "" {
		result = Flip(result, req.Flip)
	}
	if req.Filters != nil {
		if req.Filters.Grayscale {
			result = Grayscale(result)
		}
		if req.Filters.Sepia {
			result = Sepia(result)
		}
		if req.Filters.Blur > 0 {
			result = Blur(result, req.Filters.Blur)
		}
		if req.Filters.Brightness != 0 {
			result = Brightness(result, req.Filters.Brightness)
		}
		if req.Filters.Contrast != 0 {
			result = Contrast(result, req.Filters.Contrast)
		}
	}
	return result, nil
}

