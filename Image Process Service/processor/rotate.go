package processor

import "image"

func Rotate(src image.Image, angle int) image.Image {
	angle = angle % 360
	if angle < 0 {
		angle += 360
	}
	if angle == 0 {
		return src
	}
	bounds := src.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	switch angle {
	case 90:
		dst := image.NewRGBA(image.Rect(0, 0, height, width))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				dst.Set(height-1-y, x, src.At(bounds.Min.X+x, bounds.Min.Y+y))
			}
		}
		return dst
	case 180:
		dst := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				dst.Set(width-1-x, height-1-y, src.At(bounds.Min.X+x, bounds.Min.Y+y))
			}
		}
		return dst
	case 270:
		dst := image.NewRGBA(image.Rect(0, 0, height, width))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				dst.Set(y, width-1-x, src.At(bounds.Min.X+x, bounds.Min.Y+y))
			}
		}
		return dst
	default:
		return src
	}
}
