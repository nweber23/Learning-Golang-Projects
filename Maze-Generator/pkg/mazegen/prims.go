package mazegen

import "math/rand"

// PrimsGenerator implements a randomised Prim's algorithm.
// It produces mazes with more branches and shorter dead ends than DFS.
type PrimsGenerator struct {
	opts Options
}

func (g *PrimsGenerator) Generate(m *Maze) error {
	rng := rand.New(rand.NewSource(g.opts.Seed))

	// Pattern42Sentinel (0xFF) is defined in maze.go — Prim's reads it here.
	inMaze := make([][]bool, m.Height)
	for y := range inMaze {
		inMaze[y] = make([]bool, m.Width)
	}

	// Mark pattern42 cells as already in-maze (unbreakable) and restore sentinel to 0xF.
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.Cells[y][x] == Pattern42Sentinel {
				inMaze[y][x] = true
				m.Cells[y][x] = WallNorth | WallEast | WallSouth | WallWest
			}
		}
	}

	type edge struct {
		from Point
		dir  Direction
	}

	addFrontier := func(p Point, frontier *[]edge) {
		for _, dir := range Directions {
			delta := DirDelta[dir]
			q := Point{p.X + delta.X, p.Y + delta.Y}
			if m.InBounds(q) && !inMaze[q.Y][q.X] {
				*frontier = append(*frontier, edge{p, dir})
			}
		}
	}

	// Start from Entry.
	inMaze[m.Entry.Y][m.Entry.X] = true
	var frontier []edge
	addFrontier(m.Entry, &frontier)

	for len(frontier) > 0 {
		// Pick a random frontier edge.
		idx := rng.Intn(len(frontier))
		e := frontier[idx]
		// Remove by swapping with last element.
		frontier[idx] = frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]

		delta := DirDelta[e.dir]
		q := Point{e.from.X + delta.X, e.from.Y + delta.Y}

		if inMaze[q.Y][q.X] {
			continue // already connected
		}

		m.RemoveWall(e.from, e.dir)
		inMaze[q.Y][q.X] = true
		addFrontier(q, &frontier)
	}

	if !g.opts.Perfect {
		// Reuse DFS extra-passage logic.
		dfsGen := &DFSGenerator{opts: g.opts}
		dfsGen.addExtraPassages(m, rng)
	}

	return nil
}
