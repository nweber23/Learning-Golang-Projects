package render

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"

	"github.com/niklaswebde/maze-generator/pkg/mazegen"
)

// ANSI color codes for wall color cycling.
var wallColors = []string{
	"\033[37m", // white
	"\033[36m", // cyan
	"\033[33m", // yellow
	"\033[35m", // magenta
	"\033[32m", // green
	"\033[34m", // blue
}

const (
	ansiReset   = "\033[0m"
	clearScreen = "\033[2J\033[H"
	hideCursor  = "\033[?25l"
	showCursor  = "\033[?25h"
)

// TerminalState holds the interactive renderer state.
type TerminalState struct {
	ColorIdx int
	ShowPath bool
}

// cellWidth selects compact (2 chars/cell) or normal (3 chars/cell) mode
// based on whether the maze fits in the terminal width.
// Returns 2 or 3.
func cellWidth(mazeWidth, termWidth int) int {
	// Full width: mazeWidth * 3 (content) + mazeWidth + 1 (walls) = mazeWidth*4+1
	if mazeWidth*4+1 <= termWidth {
		return 3
	}
	// Compact width: mazeWidth * 1 (content) + mazeWidth + 1 (walls) = mazeWidth*2+1
	return 1
}

// renderWidth returns the total character width of the maze at a given cell width.
func renderWidth(mazeWidth, cw int) int {
	return mazeWidth*(cw+1) + 1
}

// renderHeight returns the total line count of the maze render.
func renderHeight(mazeHeight int) int {
	// 1 top border + mazeHeight content rows + (mazeHeight-1) separator rows + 1 bottom border + 2 legend lines
	return 1 + mazeHeight*2 + 2
}

// RenderMaze draws the maze into a bytes.Buffer and writes it atomically to stdout.
// It auto-selects compact vs normal cell width based on terminal size.
func RenderMaze(m *mazegen.Maze, state TerminalState) {
	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || termWidth == 0 {
		termWidth, termHeight = 80, 24 // safe fallback
	}

	cw := cellWidth(m.Width, termWidth)
	rw := renderWidth(m.Width, cw)
	rh := renderHeight(m.Height)

	// Horizontal padding to centre the maze.
	hpad := 0
	if rw < termWidth {
		hpad = (termWidth - rw) / 2
	}
	padStr := ""
	for i := 0; i < hpad; i++ {
		padStr += " "
	}

	// Vertical padding to centre the maze.
	vpad := 0
	if rh < termHeight {
		vpad = (termHeight - rh) / 2
	}

	color := wallColors[state.ColorIdx%len(wallColors)]

	// Solution set for fast O(1) lookup.
	pathSet := make(map[mazegen.Point]bool)
	if state.ShowPath && m.Solution != nil {
		cur := m.Entry
		pathSet[cur] = true
		for _, dir := range m.Solution {
			d := mazegen.DirDelta[dir]
			cur = mazegen.Point{X: cur.X + d.X, Y: cur.Y + d.Y}
			pathSet[cur] = true
		}
	}

	var buf bytes.Buffer

	// Clear screen and move to top-left.
	buf.WriteString(clearScreen)

	// If maze is too small to display at all, show an error.
	if rw > termWidth || rh > termHeight {
		msg := fmt.Sprintf(" Maze too large for terminal (%dx%d). Resize or use a smaller maze.", m.Width, m.Height)
		for i := 0; i < termHeight/2; i++ {
			buf.WriteString("\n")
		}
		buf.WriteString(msg)
		buf.WriteString("\n\n [q] quit  [r] regenerate")
		os.Stdout.Write(buf.Bytes())
		return
	}

	// Vertical top padding.
	for i := 0; i < vpad; i++ {
		buf.WriteString("\n")
	}

	// Helper: write a padded line.
	line := func(s string) {
		buf.WriteString(padStr)
		buf.WriteString(s)
		buf.WriteString("\n")
	}

	// wall/space segments for horizontal runs.
	wallSeg := func() string {
		s := ""
		for i := 0; i < cw; i++ {
			s += "─"
		}
		return s
	}
	spaceSeg := func() string {
		s := ""
		for i := 0; i < cw; i++ {
			s += " "
		}
		return s
	}
	ws := wallSeg()
	ss := spaceSeg()

	// Cell content string.
	cellContent := func(p mazegen.Point) string {
		switch {
		case p == m.Entry:
			if cw == 1 {
				return "E"
			}
			return " E "
		case p == m.Exit:
			if cw == 1 {
				return "X"
			}
			return " X "
		case pathSet[p]:
			if cw == 1 {
				return "\033[31m▪\033[0m"
			}
			return "\033[31m ▪ \033[0m"
		default:
			if cw == 1 {
				return " "
			}
			return " · "
		}
	}

	// ── Top border ──
	{
		row := color + "┌"
		for x := 0; x < m.Width; x++ {
			if m.HasWall(mazegen.Point{X: x, Y: 0}, mazegen.North) {
				row += ws
			} else {
				row += ss
			}
			if x < m.Width-1 {
				row += "┬"
			}
		}
		row += "┐" + ansiReset
		line(row)
	}

	for y := 0; y < m.Height; y++ {
		// ── Cell content row ──
		{
			row := color + "│" + ansiReset
			for x := 0; x < m.Width; x++ {
				p := mazegen.Point{X: x, Y: y}
				row += cellContent(p)
				if m.HasWall(p, mazegen.East) {
					row += color + "│" + ansiReset
				} else {
					row += " "
				}
			}
			line(row)
		}

		// ── Horizontal separator row ──
		if y < m.Height-1 {
			row := color + "├"
			for x := 0; x < m.Width; x++ {
				if m.HasWall(mazegen.Point{X: x, Y: y}, mazegen.South) {
					row += ws
				} else {
					row += ss
				}
				if x < m.Width-1 {
					row += "┼"
				}
			}
			row += "┤" + ansiReset
			line(row)
		}
	}

	// ── Bottom border ──
	{
		row := color + "└"
		for x := 0; x < m.Width; x++ {
			if m.HasWall(mazegen.Point{X: x, Y: m.Height - 1}, mazegen.South) {
				row += ws
			} else {
				row += ss
			}
			if x < m.Width-1 {
				row += "┴"
			}
		}
		row += "┘" + ansiReset
		line(row)
	}

	// ── Key legend ──
	buf.WriteString("\n")
	line(" [r] regenerate  [p] path  [c] color  [q] quit")

	// Write everything atomically — no flickering.
	os.Stdout.Write(buf.Bytes())
}

// RunInteractive enters raw terminal mode and handles keypresses until 'q'.
// onRegen is called to produce a new maze.
// onQuit is called when the user presses 'q'.
func RunInteractive(m *mazegen.Maze, state TerminalState, onRegen func() *mazegen.Maze, onQuit func()) {
	fd := int(os.Stdin.Fd())
	old, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not enter raw terminal mode: %v\n", err)
		return
	}
	defer func() {
		term.Restore(fd, old)
		fmt.Print(showCursor)
	}()

	fmt.Print(hideCursor)
	RenderMaze(m, state)

	// Listen for terminal resize (SIGWINCH) and re-render.
	sigwinch := make(chan os.Signal, 1)
	signal.Notify(sigwinch, syscall.SIGWINCH)
	defer signal.Stop(sigwinch)

	go func() {
		for range sigwinch {
			RenderMaze(m, state)
		}
	}()

	buf := make([]byte, 3) // 3 bytes to catch escape sequences
	for {
		n, _ := os.Stdin.Read(buf)
		if n == 0 {
			continue
		}
		switch buf[0] {
		case 'r':
			m = onRegen()
			m.Solution = nil
			state.ShowPath = false
			RenderMaze(m, state)
		case 'p':
			state.ShowPath = !state.ShowPath
			if state.ShowPath && m.Solution == nil {
				if solveErr := m.Solve(); solveErr != nil {
					// Show briefly then re-render — no permanent stderr pollution in raw mode.
					_ = solveErr
				}
			}
			RenderMaze(m, state)
		case 'c':
			state.ColorIdx++
			RenderMaze(m, state)
		case 'q', 3: // 'q' or Ctrl-C
			onQuit()
			return
		}
	}
}
