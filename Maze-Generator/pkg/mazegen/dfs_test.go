package mazegen

import "testing"

func TestDFS_Coherence(t *testing.T) {
	m := makeTestMaze(15, 10)
	g := &DFSGenerator{opts: Options{Perfect: true, Seed: 1}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	checkCoherence(t, m)
}

func TestDFS_Connectivity(t *testing.T) {
	m := makeTestMaze(15, 10)
	g := &DFSGenerator{opts: Options{Perfect: true, Seed: 2}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	checkConnectivity(t, m)
}

func TestDFS_PerfectMaze(t *testing.T) {
	m := makeTestMaze(15, 10)
	g := &DFSGenerator{opts: Options{Perfect: true, Seed: 3}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	checkPerfect(t, m)
}

func TestDFS_ImperfectHasMorePassages(t *testing.T) {
	m := makeTestMaze(15, 10)
	g := &DFSGenerator{opts: Options{Perfect: false, Seed: 4}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	checkCoherence(t, m)
	checkConnectivity(t, m)
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
	minimum := m.Width*m.Height - 1
	if passages <= minimum {
		t.Errorf("imperfect maze should have more than %d passages, got %d", minimum, passages)
	}
}

func TestDFS_Reproducible(t *testing.T) {
	m1 := makeTestMaze(10, 10)
	m2 := makeTestMaze(10, 10)
	g1 := &DFSGenerator{opts: Options{Perfect: true, Seed: 99}}
	g2 := &DFSGenerator{opts: Options{Perfect: true, Seed: 99}}
	if err := g1.Generate(m1); err != nil {
		t.Fatal(err)
	}
	if err := g2.Generate(m2); err != nil {
		t.Fatal(err)
	}
	for y := 0; y < m1.Height; y++ {
		for x := 0; x < m1.Width; x++ {
			if m1.Cells[y][x] != m2.Cells[y][x] {
				t.Errorf("cell (%d,%d) differs: %d vs %d", x, y, m1.Cells[y][x], m2.Cells[y][x])
			}
		}
	}
}
