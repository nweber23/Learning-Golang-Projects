package mazegen

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// dirLetter maps Direction to its single-char export string.
var dirLetter = [4]string{
	North: "N",
	East:  "E",
	South: "S",
	West:  "W",
}

// letterDir maps export char to Direction.
var letterDir = map[byte]Direction{
	'N': North,
	'E': East,
	'S': South,
	'W': West,
}

// Write writes the maze to a file in the spec hex format:
//
//	One hex char per cell, row by row (outer=row, inner=col)
//	Blank line
//	Entry x,y
//	Exit x,y
//	Solution path as NESW string (empty line if Solution is nil)
func Write(m *Maze, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	// Cell rows.
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			fmt.Fprintf(w, "%X", m.Cells[y][x])
		}
		fmt.Fprintln(w)
	}

	// Blank separator.
	fmt.Fprintln(w)

	// Entry, exit.
	fmt.Fprintf(w, "%d,%d\n", m.Entry.X, m.Entry.Y)
	fmt.Fprintf(w, "%d,%d\n", m.Exit.X, m.Exit.Y)

	// Solution path (empty line if nil).
	if m.Solution != nil {
		var sb strings.Builder
		for _, dir := range m.Solution {
			sb.WriteString(dirLetter[dir])
		}
		fmt.Fprintln(w, sb.String())
	} else {
		fmt.Fprintln(w)
	}

	return w.Flush()
}

// Parse reads a hex maze file written by Write and reconstructs the Maze.
// Solution is nil if the solution line was empty.
func Parse(path string) (*Maze, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var rowLines []string
	// Read cell rows until blank line.
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		rowLines = append(rowLines, line)
	}
	if len(rowLines) == 0 {
		return nil, fmt.Errorf("no cell data in file")
	}

	height := len(rowLines)
	width := len(rowLines[0])
	cells := make([][]Cell, height)
	for y, row := range rowLines {
		if len(row) != width {
			return nil, fmt.Errorf("row %d has length %d, expected %d", y, len(row), width)
		}
		cells[y] = make([]Cell, width)
		for x, ch := range row {
			val, err := strconv.ParseUint(string(ch), 16, 8)
			if err != nil {
				return nil, fmt.Errorf("invalid hex char %q at row %d col %d", ch, y, x)
			}
			cells[y][x] = Cell(val)
		}
	}

	// Read entry.
	if !scanner.Scan() {
		return nil, fmt.Errorf("missing entry line")
	}
	entry, err := parsePoint(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("parse entry: %w", err)
	}

	// Read exit.
	if !scanner.Scan() {
		return nil, fmt.Errorf("missing exit line")
	}
	exit, err := parsePoint(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("parse exit: %w", err)
	}

	// Read solution (may be empty).
	var solution []Direction
	if scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			solution = make([]Direction, len(line))
			for i := 0; i < len(line); i++ {
				dir, ok := letterDir[line[i]]
				if !ok {
					return nil, fmt.Errorf("invalid direction char %q in solution", line[i])
				}
				solution[i] = dir
			}
		}
	}

	return &Maze{
		Width:    width,
		Height:   height,
		Cells:    cells,
		Entry:    entry,
		Exit:     exit,
		Solution: solution,
	}, nil
}

func parsePoint(s string) (Point, error) {
	parts := strings.SplitN(s, ",", 2)
	if len(parts) != 2 {
		return Point{}, fmt.Errorf("expected x,y got %q", s)
	}
	x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return Point{}, fmt.Errorf("invalid x: %w", err)
	}
	y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Point{}, fmt.Errorf("invalid y: %w", err)
	}
	return Point{x, y}, nil
}
