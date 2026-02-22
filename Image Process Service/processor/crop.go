package processor

import "image"

func Crop(src image.Image, x, y, width, height int) image.Image {
	bounds := src.Bounds()
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x+width > bounds.Max.X {
		width = bounds.Max.X - x
	}
	if y+height > bounds.Max.Y {
		height = bounds.Max.Y - y
	}
	if width <= 0 || height <= 0 {
		return src
	}
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	for cy := 0; cy < height; cy++ {
		for cx := 0; cx < width; cx++ {
			dst.Set(cx, cy, src.At(x+cx, y+cy))
		}
	}
	return dst
}
