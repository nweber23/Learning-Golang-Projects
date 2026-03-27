package mazegen

// NOTE: Pattern42Sentinel (0xFF) is defined in maze.go, not here,
// so DFS and Prim's can import it without a circular dependency.

// Pattern42MinWidth and Pattern42MinHeight are the smallest maze dimensions
// that can fit the "42" pattern.
const (
	Pattern42MinWidth  = 11
	Pattern42MinHeight = 7
)

// pattern42Pixels defines the "42" glyph as pixel offsets (col, row) relative to origin (0,0).
// Each pixel becomes a fully-closed cell. The glyph is 7 cols × 5 rows.
//
// "4" occupies cols 0-2, "2" occupies cols 4-6 (1-col gap between them).
//
//	"4":          "2":
//	# . #         # # #
//	# . #         . . #
//	# # #         # # #
//	. . #         # . .
//	. . #         # # #
var pattern42Pixels = []Point{
	// "4" digit
	{0, 0}, {2, 0},
	{0, 1}, {2, 1},
	{0, 2}, {1, 2}, {2, 2},
	{2, 3},
	{2, 4},
	// "2" digit
	{4, 0}, {5, 0}, {6, 0},
	{6, 1},
	{4, 2}, {5, 2}, {6, 2},
	{4, 3},
	{4, 4}, {5, 4}, {6, 4},
}

// pattern42Cells returns the absolute maze coordinates of all "42" cells,
// centred within the maze. Returns nil if the maze is too small.
func pattern42Cells(width, height int) []Point {
	if width < Pattern42MinWidth || height < Pattern42MinHeight {
		return nil
	}
	// Centre the 7×5 glyph within the maze.
	glyphW, glyphH := 7, 5
	startX := (width - glyphW) / 2
	startY := (height - glyphH) / 2
	cells := make([]Point, len(pattern42Pixels))
	for i, px := range pattern42Pixels {
		cells[i] = Point{startX + px.X, startY + px.Y}
	}
	return cells
}

// StampPattern42 pre-places the "42" glyph onto the maze grid before generation.
// Each glyph cell is set to Pattern42Sentinel so generators treat it as unbreakable.
// Returns (true, nil) if stamped, (false, nil) if maze is too small (caller should print warning).
func StampPattern42(m *Maze) (bool, error) {
	cells := pattern42Cells(m.Width, m.Height)
	if cells == nil {
		return false, nil
	}
	for _, p := range cells {
		m.Cells[p.Y][p.X] = Pattern42Sentinel
	}
	return true, nil
}
