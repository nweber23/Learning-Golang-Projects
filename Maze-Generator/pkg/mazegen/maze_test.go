package mazegen

import "testing"

// checkCoherence verifies that neighbouring cells agree on shared walls.
func checkCoherence(t *testing.T, m *Maze) {
	t.Helper()
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			for _, dir := range Directions {
				delta := DirDelta[dir]
				q := Point{x + delta.X, y + delta.Y}
				if !m.InBounds(q) {
					continue
				}
				hasWallHere := m.Cells[y][x]&DirWall[dir] != 0
				hasWallThere := m.Cells[q.Y][q.X]&Opposite(DirWall[dir]) != 0
				if hasWallHere != hasWallThere {
					t.Errorf("incoherent wall between (%d,%d) dir %d and (%d,%d): %v vs %v",
						x, y, dir, q.X, q.Y, hasWallHere, hasWallThere)
				}
			}
		}
	}
}

// checkConnectivity verifies all cells are reachable from the entry via BFS.
// NOTE: This helper is designed for mazes WITHOUT pattern42 stamped.
// Pattern42 cells are fully walled and will fail this check.
func checkConnectivity(t *testing.T, m *Maze) {
	t.Helper()
	visited := make([][]bool, m.Height)
	for y := range visited {
		visited[y] = make([]bool, m.Width)
	}
	queue := []Point{m.Entry}
	visited[m.Entry.Y][m.Entry.X] = true
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		for _, dir := range Directions {
			if m.HasWall(p, dir) {
				continue
			}
			delta := DirDelta[dir]
			q := Point{p.X + delta.X, p.Y + delta.Y}
			if m.InBounds(q) && !visited[q.Y][q.X] {
				visited[q.Y][q.X] = true
				queue = append(queue, q)
			}
		}
	}
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if !visited[y][x] {
				t.Errorf("cell (%d,%d) is not reachable from entry %v", x, y, m.Entry)
			}
		}
	}
}

// checkPerfect verifies the maze has exactly (width*height - 1) passages (spanning tree).
// NOTE: Use only on mazes WITHOUT pattern42 stamped.
func checkPerfect(t *testing.T, m *Maze) {
	t.Helper()
	passages := 0
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			p := Point{x, y}
			for _, dir := range []Direction{East, South} {
				if !m.HasWall(p, dir) {
					q := Point{x + DirDelta[dir].X, y + DirDelta[dir].Y}
					if m.InBounds(q) {
						passages++
					}
				}
			}
		}
	}
	expected := m.Width*m.Height - 1
	if passages != expected {
		t.Errorf("perfect maze should have %d passages, got %d", expected, passages)
	}
}

// makeTestMaze creates a simple maze for testing (no pattern42).
func makeTestMaze(w, h int) *Maze {
	return NewMaze(w, h, Point{0, 0}, Point{w - 1, h - 1}, 42)
}
