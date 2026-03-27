package render

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/niklaswebde/maze-generator/pkg/mazegen"
)

const (
	cellSize  = 16 // pixels per cell
	wallWidth = 2  // pixels for wall lines
)

// WritePNG renders the maze to a PNG file at the given path.
// If m.Solution is non-nil, the solution path is highlighted in red.
func WritePNG(m *mazegen.Maze, path string) error {
	imgW := m.Width*cellSize + wallWidth
	imgH := m.Height*cellSize + wallWidth

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))

	// Fill background white.
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	wallColor := color.RGBA{R: 30, G: 30, B: 30, A: 255}
	pathColor := color.RGBA{R: 220, G: 60, B: 60, A: 255}
	entryColor := color.RGBA{R: 60, G: 180, B: 60, A: 255}
	exitColor := color.RGBA{R: 60, G: 60, B: 220, A: 255}

	// Draw walls for each cell.
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			px := x*cellSize + wallWidth
			py := y*cellSize + wallWidth
			p := mazegen.Point{X: x, Y: y}

			if m.HasWall(p, mazegen.North) {
				fillRect(img, px-wallWidth, py-wallWidth, cellSize+wallWidth, wallWidth, wallColor)
			}
			if m.HasWall(p, mazegen.West) {
				fillRect(img, px-wallWidth, py-wallWidth, wallWidth, cellSize+wallWidth, wallColor)
			}
			// Always draw south border for last row.
			if y == m.Height-1 && m.HasWall(p, mazegen.South) {
				fillRect(img, px-wallWidth, py+cellSize-wallWidth, cellSize+wallWidth, wallWidth, wallColor)
			}
			// Always draw east border for last col.
			if x == m.Width-1 && m.HasWall(p, mazegen.East) {
				fillRect(img, px+cellSize-wallWidth, py-wallWidth, wallWidth, cellSize+wallWidth, wallColor)
			}
		}
	}

	// Draw solution path.
	if m.Solution != nil {
		cur := m.Entry
		for _, dir := range m.Solution {
			cx := cur.X*cellSize + wallWidth + cellSize/2
			cy := cur.Y*cellSize + wallWidth + cellSize/2
			delta := mazegen.DirDelta[dir]
			next := mazegen.Point{X: cur.X + delta.X, Y: cur.Y + delta.Y}
			nx := next.X*cellSize + wallWidth + cellSize/2
			ny := next.Y*cellSize + wallWidth + cellSize/2
			drawLine(img, cx, cy, nx, ny, pathColor)
			cur = next
		}
	}

	// Mark entry (green) and exit (blue).
	markCell(img, m.Entry, entryColor)
	markCell(img, m.Exit, exitColor)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func fillRect(img *image.RGBA, x, y, w, h int, c color.Color) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, c)
		}
	}
}

func markCell(img *image.RGBA, p mazegen.Point, c color.Color) {
	cx := p.X*cellSize + wallWidth + cellSize/4
	cy := p.Y*cellSize + wallWidth + cellSize/4
	fillRect(img, cx, cy, cellSize/2, cellSize/2, c)
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	dx := absInt(x1 - x0)
	dy := absInt(y1 - y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy
	for {
		img.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
