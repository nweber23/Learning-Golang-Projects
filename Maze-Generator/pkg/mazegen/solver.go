package mazegen

import "fmt"

// Solve finds the shortest path from Entry to Exit using BFS and stores it in m.Solution.
// Calling Solve multiple times is safe; it recomputes and overwrites m.Solution.
// Returns an error if no path exists (should not happen in a valid fully-connected maze).
func (m *Maze) Solve() error {
	type state struct {
		p    Point
		from *state
		dir  Direction
	}

	visited := make([][]bool, m.Height)
	for y := range visited {
		visited[y] = make([]bool, m.Width)
	}

	start := &state{p: m.Entry}
	queue := []*state{start}
	visited[m.Entry.Y][m.Entry.X] = true

	var found *state
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.p == m.Exit {
			found = cur
			break
		}
		for _, dir := range Directions {
			if m.HasWall(cur.p, dir) {
				continue
			}
			delta := DirDelta[dir]
			q := Point{cur.p.X + delta.X, cur.p.Y + delta.Y}
			if m.InBounds(q) && !visited[q.Y][q.X] {
				visited[q.Y][q.X] = true
				queue = append(queue, &state{p: q, from: cur, dir: dir})
			}
		}
	}

	if found == nil {
		return fmt.Errorf("no path from entry %v to exit %v", m.Entry, m.Exit)
	}

	// Reconstruct path by walking back through `from` pointers.
	var path []Direction
	for s := found; s.from != nil; s = s.from {
		path = append(path, s.dir)
	}
	// Reverse.
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	m.Solution = path
	return nil
}
