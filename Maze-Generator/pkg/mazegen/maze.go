// Package mazegen provides types and algorithms for generating and solving mazes.
// The core Maze type uses a bitmask-per-cell wall representation compatible with
// the A-Maze-ing hex output format.
package mazegen

// Cell encodes which walls are closed as a 4-bit bitmask.
// Bit 0 (LSB) = North, Bit 1 = East, Bit 2 = South, Bit 3 = West.
// A set bit means the wall is CLOSED (present).
type Cell uint8

const (
	WallNorth Cell = 1 << iota // 0001
	WallEast                   // 0010
	WallSouth                  // 0100
	WallWest                   // 1000
)

// Opposite returns the wall on the opposite side of the given wall.
func Opposite(wall Cell) Cell {
	switch wall {
	case WallNorth:
		return WallSouth
	case WallSouth:
		return WallNorth
	case WallEast:
		return WallWest
	case WallWest:
		return WallEast
	}
	return 0
}

// Point is a 2D grid coordinate (column, row).
// Defined here to keep pkg/mazegen self-contained (no image package import).
type Point struct{ X, Y int }

// Direction represents a cardinal movement direction.
type Direction uint8

const (
	North Direction = iota
	East
	South
	West
)

// DirDelta maps a Direction to its (dX, dY) grid delta.
// Fixed array — index by Direction constant. Do not modify.
var DirDelta = [4]Point{
	North: {0, -1},
	East:  {1, 0},
	South: {0, 1},
	West:  {-1, 0},
}

// DirWall maps a Direction to the wall bit on the current cell.
// Fixed array — index by Direction constant. Do not modify.
var DirWall = [4]Cell{
	North: WallNorth,
	East:  WallEast,
	South: WallSouth,
	West:  WallWest,
}

// Directions is the canonical list of all cardinal directions.
// Callers must copy before shuffling — do not modify this slice.
var Directions = [4]Direction{North, East, South, West}

// Pattern42Sentinel is written to Cells[y][x] to mark a pre-placed "42" cell.
// Defined here (not in pattern42.go) so DFS and Prim's can read it without a circular import.
// 0xFF is outside the 4-bit wall range (0x0–0xF) and is never a valid generated wall value.
const Pattern42Sentinel Cell = 0xFF

// Maze represents a generated maze grid.
// Cells is row-major: Cells[row][col], i.e. Cells[y][x].
// Solution is nil until Solve() is called explicitly.
type Maze struct {
	Width, Height int
	Cells         [][]Cell
	Entry, Exit   Point
	Seed          int64
	Solution      []Direction
}

// NewMaze allocates a Maze with all walls closed (all cells = 0xF).
func NewMaze(width, height int, entry, exit Point, seed int64) *Maze {
	cells := make([][]Cell, height)
	for y := range cells {
		cells[y] = make([]Cell, width)
		for x := range cells[y] {
			cells[y][x] = WallNorth | WallEast | WallSouth | WallWest
		}
	}
	return &Maze{
		Width:  width,
		Height: height,
		Cells:  cells,
		Entry:  entry,
		Exit:   exit,
		Seed:   seed,
	}
}

// InBounds reports whether p is within the maze grid.
func (m *Maze) InBounds(p Point) bool {
	return p.X >= 0 && p.X < m.Width && p.Y >= 0 && p.Y < m.Height
}

// RemoveWall removes the wall between two adjacent cells (both sides).
// p and q must be adjacent. Does nothing if either is out of bounds.
func (m *Maze) RemoveWall(p Point, dir Direction) {
	delta := DirDelta[dir]
	q := Point{p.X + delta.X, p.Y + delta.Y}
	if !m.InBounds(p) || !m.InBounds(q) {
		return
	}
	m.Cells[p.Y][p.X] &^= DirWall[dir]
	m.Cells[q.Y][q.X] &^= Opposite(DirWall[dir])
}

// HasWall reports whether the wall in the given direction is closed.
func (m *Maze) HasWall(p Point, dir Direction) bool {
	if !m.InBounds(p) {
		return true // out-of-bounds treated as solid wall
	}
	return m.Cells[p.Y][p.X]&DirWall[dir] != 0
}
