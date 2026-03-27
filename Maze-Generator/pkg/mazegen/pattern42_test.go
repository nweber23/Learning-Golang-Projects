package mazegen

import "testing"

func TestPattern42_CellsFullyClosed(t *testing.T) {
	m := makeTestMaze(20, 12)
	stamped, err := StampPattern42(m)
	if err != nil {
		t.Fatal(err)
	}
	if !stamped {
		t.Fatal("expected pattern to be stamped on 20x12 maze")
	}
	// After DFS generation, all "42" cells must still have all 4 walls closed.
	g := &DFSGenerator{opts: Options{Perfect: true, Seed: 1}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	for _, p := range pattern42Cells(m.Width, m.Height) {
		if m.Cells[p.Y][p.X] != (WallNorth | WallEast | WallSouth | WallWest) {
			t.Errorf("pattern42 cell (%d,%d) should be fully closed, got %04b",
				p.X, p.Y, m.Cells[p.Y][p.X])
		}
	}
}

func TestPattern42_TooSmall(t *testing.T) {
	m := makeTestMaze(5, 5)
	stamped, err := StampPattern42(m)
	if err != nil {
		t.Fatal(err)
	}
	if stamped {
		t.Error("expected no stamp on 5x5 maze (too small)")
	}
}

func TestPattern42_WithinBounds(t *testing.T) {
	m := makeTestMaze(15, 9)
	StampPattern42(m)
	for _, p := range pattern42Cells(m.Width, m.Height) {
		if !m.InBounds(p) {
			t.Errorf("pattern42 cell (%d,%d) is out of bounds for %dx%d maze",
				p.X, p.Y, m.Width, m.Height)
		}
	}
}
