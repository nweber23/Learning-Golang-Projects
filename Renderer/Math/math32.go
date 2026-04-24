package Math

import (
	"math"
)

const (
	pi32 = float32(math.Pi)
)

func sqrt32(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func sin32(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func cos32(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func tan32(x float32) float32 {
	return float32(math.Tan(float64(x)))
}
