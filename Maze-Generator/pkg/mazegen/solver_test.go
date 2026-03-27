package mazegen

import "testing"

func TestSolver_FindsPath(t *testing.T) {
	m := makeTestMaze(10, 10)
	g := &DFSGenerator{opts: Options{Perfect: true, Seed: 5}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	if err := m.Solve(); err != nil {
		t.Fatalf("Solve failed: %v", err)
	}
	if m.Solution == nil {
		t.Fatal("Solution is nil after Solve()")
	}
	if len(m.Solution) == 0 {
		t.Fatal("Solution is empty")
	}
}

func TestSolver_PathIsValid(t *testing.T) {
	m := makeTestMaze(10, 10)
	g := &DFSGenerator{opts: Options{Perfect: true, Seed: 6}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	if err := m.Solve(); err != nil {
		t.Fatal(err)
	}
	cur := m.Entry
	for i, dir := range m.Solution {
		if m.HasWall(cur, dir) {
			t.Errorf("step %d: wall exists in direction %d at (%d,%d)", i, dir, cur.X, cur.Y)
		}
		delta := DirDelta[dir]
		cur = Point{cur.X + delta.X, cur.Y + delta.Y}
		if !m.InBounds(cur) {
			t.Fatalf("step %d: moved out of bounds to (%d,%d)", i, cur.X, cur.Y)
		}
	}
	if cur != m.Exit {
		t.Errorf("path ends at (%d,%d), expected Exit (%d,%d)", cur.X, cur.Y, m.Exit.X, m.Exit.Y)
	}
}

func TestSolver_ShortestPath(t *testing.T) {
	m := makeTestMaze(8, 8)
	g := &DFSGenerator{opts: Options{Perfect: true, Seed: 7}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	if err := m.Solve(); err != nil {
		t.Fatal(err)
	}
	// Independent BFS to verify path length.
	type state struct {
		p    Point
		dist int
	}
	visited := make([][]bool, m.Height)
	for y := range visited {
		visited[y] = make([]bool, m.Width)
	}
	queue := []state{{m.Entry, 0}}
	visited[m.Entry.Y][m.Entry.X] = true
	bfsDist := -1
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		if s.p == m.Exit {
			bfsDist = s.dist
			break
		}
		for _, dir := range Directions {
			if m.HasWall(s.p, dir) {
				continue
			}
			delta := DirDelta[dir]
			q := Point{s.p.X + delta.X, s.p.Y + delta.Y}
			if m.InBounds(q) && !visited[q.Y][q.X] {
				visited[q.Y][q.X] = true
				queue = append(queue, state{q, s.dist + 1})
			}
		}
	}
	if len(m.Solution) != bfsDist {
		t.Errorf("solution length %d != BFS distance %d", len(m.Solution), bfsDist)
	}
}

func TestSolver_Idempotent(t *testing.T) {
	m := makeTestMaze(8, 8)
	g := &DFSGenerator{opts: Options{Perfect: true, Seed: 8}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	if err := m.Solve(); err != nil {
		t.Fatal(err)
	}
	sol1 := make([]Direction, len(m.Solution))
	copy(sol1, m.Solution)
	if err := m.Solve(); err != nil {
		t.Fatal(err)
	}
	if len(sol1) != len(m.Solution) {
		t.Errorf("second Solve() gave different length: %d vs %d", len(sol1), len(m.Solution))
	}
}
