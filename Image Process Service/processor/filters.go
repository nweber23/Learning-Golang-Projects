package processor

import (
	"image"
	"image/color"
	"math"
)

func Grayscale(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			gray := (r*299 + g*587 + b*114) / 1000

			dst.Set(x, y, color.RGBA64{R: uint16(gray), G: uint16(gray), B: uint16(gray), A: uint16(a)})
		}
	}

	return dst
}

func Sepia(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()

			outR := uint16(math.Min(65535, float64(r)*0.393+float64(g)*0.769+float64(b)*0.189))
			outG := uint16(math.Min(65535, float64(r)*0.349+float64(g)*0.686+float64(b)*0.168))
			outB := uint16(math.Min(65535, float64(r)*0.272+float64(g)*0.534+float64(b)*0.131))

			dst.Set(x, y, color.RGBA64{R: outR, G: outG, B: outB, A: uint16(a)})
		}
	}

	return dst
}

func Blur(src image.Image, radius float64) image.Image {
	if radius <= 0 {
		return src
	}

	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	r := int(radius)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var sumR, sumG, sumB, sumA float64
			var count int

			for ky := -r; ky <= r; ky++ {
				for kx := -r; kx <= r; kx++ {
					px := x + kx
					py := y + ky

					if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
						cr, cg, cb, ca := src.At(px, py).RGBA()
						sumR += float64(cr)
						sumG += float64(cg)
						sumB += float64(cb)
						sumA += float64(ca)
						count++
					}
				}
			}

			dst.Set(x, y, color.RGBA64{
				R: uint16(sumR / float64(count)),
				G: uint16(sumG / float64(count)),
				B: uint16(sumB / float64(count)),
				A: uint16(sumA / float64(count)),
			})
		}
	}

	return dst
}

func Brightness(src image.Image, factor float64) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()

			outR := uint16(math.Max(0, math.Min(65535, float64(r)*factor)))
			outG := uint16(math.Max(0, math.Min(65535, float64(g)*factor)))
			outB := uint16(math.Max(0, math.Min(65535, float64(b)*factor)))

			dst.Set(x, y, color.RGBA64{R: outR, G: outG, B: outB, A: uint16(a)})
		}
	}

	return dst
}

func Contrast(src image.Image, factor float64) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()

			outR := uint16(math.Max(0, math.Min(65535, (float64(r)-32768)*factor+32768)))
			outG := uint16(math.Max(0, math.Min(65535, (float64(g)-32768)*factor+32768)))
			outB := uint16(math.Max(0, math.Min(65535, (float64(b)-32768)*factor+32768)))

			dst.Set(x, y, color.RGBA64{R: outR, G: outG, B: outB, A: uint16(a)})
		}
	}

	return dst
}
