package mazegen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExport_RoundTrip_WithSolution(t *testing.T) {
	m := makeTestMaze(10, 8)
	g := &DFSGenerator{opts: Options{Perfect: true, Seed: 20}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	if err := m.Solve(); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.TempDir(), "maze.txt")
	if err := Write(m, path); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	m2, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if m2.Width != m.Width || m2.Height != m.Height {
		t.Errorf("dimensions mismatch: got %dx%d, want %dx%d", m2.Width, m2.Height, m.Width, m.Height)
	}
	if m2.Entry != m.Entry || m2.Exit != m.Exit {
		t.Errorf("entry/exit mismatch: got entry=%v exit=%v, want entry=%v exit=%v",
			m2.Entry, m2.Exit, m.Entry, m.Exit)
	}
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m2.Cells[y][x] != m.Cells[y][x] {
				t.Errorf("cell (%d,%d): got %X, want %X", x, y, m2.Cells[y][x], m.Cells[y][x])
			}
		}
	}
	if len(m2.Solution) != len(m.Solution) {
		t.Errorf("solution length mismatch: got %d, want %d", len(m2.Solution), len(m.Solution))
	}
}

func TestExport_RoundTrip_NilSolution(t *testing.T) {
	m := makeTestMaze(8, 6)
	g := &DFSGenerator{opts: Options{Perfect: true, Seed: 21}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}
	// Do NOT call Solve() — Solution stays nil.

	path := filepath.Join(t.TempDir(), "maze_nosol.txt")
	if err := Write(m, path); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	m2, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if m2.Solution != nil {
		t.Errorf("expected nil Solution after round-trip with no solve, got %v", m2.Solution)
	}
}

func TestExport_FileFormat(t *testing.T) {
	m := makeTestMaze(3, 2)
	g := &DFSGenerator{opts: Options{Perfect: true, Seed: 22}}
	if err := g.Generate(m); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.TempDir(), "small.txt")
	if err := Write(m, path); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := splitLines(string(data))
	if len(lines) < 6 {
		t.Errorf("expected at least 6 lines, got %d:\n%s", len(lines), string(data))
	}
	if len(lines[0]) != 3 {
		t.Errorf("row 0 should have 3 hex chars, got %q", lines[0])
	}
	if len(lines[1]) != 3 {
		t.Errorf("row 1 should have 3 hex chars, got %q", lines[1])
	}
	if lines[2] != "" {
		t.Errorf("line 3 should be blank, got %q", lines[2])
	}
}

// splitLines splits content by newline, omitting the trailing empty element.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	return lines
}
