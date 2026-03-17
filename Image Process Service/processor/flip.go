package processor

import "image"

func Flip(src image.Image, direction string) image.Image {
	bounds := src.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	dst := image.NewRGBA(bounds)
	switch direction {
	case "horizontal":
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				dst.Set(x, y, src.At(bounds.Min.X+width-1-x, bounds.Min.Y+y))
			}
		}
	case "vertical":
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				dst.Set(x, y, src.At(bounds.Min.X+x, bounds.Min.Y+height-1-y))
			}
		}
	default:
		return src
	}
	return dst
}
