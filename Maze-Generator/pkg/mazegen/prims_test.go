package mazegen

import "testing"

func TestPrims_Coherence(t *testing.T) {
	m := makeTestMaze(15, 10)
	g := &PrimsGenerator{opts: Options{Perfect: true, Seed: 10}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	checkCoherence(t, m)
}

func TestPrims_Connectivity(t *testing.T) {
	m := makeTestMaze(15, 10)
	g := &PrimsGenerator{opts: Options{Perfect: true, Seed: 11}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	checkConnectivity(t, m)
}

func TestPrims_PerfectMaze(t *testing.T) {
	m := makeTestMaze(15, 10)
	g := &PrimsGenerator{opts: Options{Perfect: true, Seed: 12}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	checkPerfect(t, m)
}

func TestPrims_Reproducible(t *testing.T) {
	m1 := makeTestMaze(10, 10)
	m2 := makeTestMaze(10, 10)
	g1 := &PrimsGenerator{opts: Options{Perfect: true, Seed: 77}}
	g2 := &PrimsGenerator{opts: Options{Perfect: true, Seed: 77}}
	if err := g1.Generate(m1); err != nil {
		t.Fatal(err)
	}
	if err := g2.Generate(m2); err != nil {
		t.Fatal(err)
	}
	for y := 0; y < m1.Height; y++ {
		for x := 0; x < m1.Width; x++ {
			if m1.Cells[y][x] != m2.Cells[y][x] {
				t.Errorf("cell (%d,%d) differs", x, y)
			}
		}
	}
}
