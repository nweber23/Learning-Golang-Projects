package processor

import "image"

func Resize(src image.Image, width, height int) image.Image {
	if width <= 0 || height <= 0 {
		return src
	}
	bounds := src.Bounds()
	srcWidth := bounds.Max.X - bounds.Min.X
	srcHeight := bounds.Max.Y - bounds.Min.Y
	if srcWidth == width && srcHeight == height {
		return src
	}
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := bounds.Min.X + (x*srcWidth)/width
			srcY := bounds.Min.Y + (y*srcHeight)/height
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}
