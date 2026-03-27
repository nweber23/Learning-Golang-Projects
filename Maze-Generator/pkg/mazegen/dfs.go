package mazegen

import (
	"math/rand"
)

// DFSGenerator implements the Recursive Backtracker (depth-first search) algorithm.
// It produces mazes with long winding corridors and high dead-end frequency.
type DFSGenerator struct {
	opts Options
}

func (g *DFSGenerator) Generate(m *Maze) error {
	rng := rand.New(rand.NewSource(g.opts.Seed))

	// Mark "42" cells (sentinel Pattern42Sentinel=0xFF, set by StampPattern42 before Generate)
	// as visited so DFS never carves through them. Restore their value to all-walls-closed (0xF)
	// so export and renderers see a normal fully-closed cell.
	visited := make([][]bool, m.Height)
	for y := range visited {
		visited[y] = make([]bool, m.Width)
	}
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.Cells[y][x] == Pattern42Sentinel {
				visited[y][x] = true
				// Restore actual wall value: all 4 walls closed.
				m.Cells[y][x] = WallNorth | WallEast | WallSouth | WallWest
			}
		}
	}

	// DFS from Entry.
	var dfs func(p Point)
	dfs = func(p Point) {
		visited[p.Y][p.X] = true
		dirs := Directions
		rng.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })
		for _, dir := range dirs {
			delta := DirDelta[dir]
			q := Point{p.X + delta.X, p.Y + delta.Y}
			if m.InBounds(q) && !visited[q.Y][q.X] {
				m.RemoveWall(p, dir)
				dfs(q)
			}
		}
	}
	dfs(m.Entry)

	// Imperfect mode: remove ~15% of remaining interior walls.
	if !g.opts.Perfect {
		g.addExtraPassages(m, rng)
	}

	return nil
}

// addExtraPassages removes approximately 15% of remaining interior walls,
// skipping any removal that would create a 3x3+ open area.
// If all candidates would create a 3x3 area, the loop exits early with fewer removals —
// this is acceptable and not an error condition.
func (g *DFSGenerator) addExtraPassages(m *Maze, rng *rand.Rand) {
	type edge struct {
		p   Point
		dir Direction
	}
	var walls []edge
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			p := Point{x, y}
			for _, dir := range []Direction{East, South} {
				if m.HasWall(p, dir) {
					delta := DirDelta[dir]
					q := Point{x + delta.X, y + delta.Y}
					if m.InBounds(q) {
						walls = append(walls, edge{p, dir})
					}
				}
			}
		}
	}
	rng.Shuffle(len(walls), func(i, j int) { walls[i], walls[j] = walls[j], walls[i] })
	target := len(walls) * 15 / 100
	removed := 0
	for _, e := range walls {
		if removed >= target {
			break
		}
		m.RemoveWall(e.p, e.dir)
		if hasOpenArea3x3(m) {
			// Undo the removal.
			delta := DirDelta[e.dir]
			q := Point{e.p.X + delta.X, e.p.Y + delta.Y}
			m.Cells[e.p.Y][e.p.X] |= DirWall[e.dir]
			m.Cells[q.Y][q.X] |= Opposite(DirWall[e.dir])
		} else {
			removed++
		}
	}
}

// hasOpenArea3x3 checks whether any 3x3 sub-grid is fully open (no interior walls).
func hasOpenArea3x3(m *Maze) bool {
	for y := 0; y <= m.Height-3; y++ {
		for x := 0; x <= m.Width-3; x++ {
			if is3x3Open(m, x, y) {
				return true
			}
		}
	}
	return false
}

// is3x3Open returns true if the 3x3 block starting at (x,y) has no internal walls.
func is3x3Open(m *Maze, x, y int) bool {
	// Check all horizontal internal passages (East walls between cols 0-1 and 1-2 for each row).
	for row := y; row < y+3; row++ {
		for col := x; col < x+2; col++ {
			if m.HasWall(Point{col, row}, East) {
				return false
			}
		}
	}
	// Check all vertical internal passages (South walls between rows 0-1 and 1-2 for each col).
	for col := x; col < x+3; col++ {
		for row := y; row < y+2; row++ {
			if m.HasWall(Point{col, row}, South) {
				return false
			}
		}
	}
	return true
}
