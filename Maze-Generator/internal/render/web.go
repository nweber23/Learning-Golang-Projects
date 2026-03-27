package render

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/niklaswebde/maze-generator/pkg/mazegen"
)

// WebServer serves an SVG representation of the maze over HTTP.
type WebServer struct {
	server *http.Server
	maze   *mazegen.Maze
}

// NewWebServer creates and starts an HTTP server on a random available port.
// Returns the server and its URL (e.g. "http://127.0.0.1:52341").
func NewWebServer(m *mazegen.Maze) (*WebServer, string, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, "", fmt.Errorf("bind web server: %w", err)
	}

	ws := &WebServer{maze: m}
	mux := http.NewServeMux()
	mux.HandleFunc("/", ws.handleMaze)
	ws.server = &http.Server{Handler: mux}

	go ws.server.Serve(ln) //nolint:errcheck

	url := fmt.Sprintf("http://%s", ln.Addr().String())
	return ws, url, nil
}

// Shutdown stops the web server gracefully.
func (ws *WebServer) Shutdown() {
	ws.server.Shutdown(context.Background()) //nolint:errcheck
}

// OpenBrowser tries to open the given URL in the default browser.
func OpenBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "linux":
		cmd, args = "xdg-open", []string{url}
	case "darwin":
		cmd, args = "open", []string{url}
	default:
		fmt.Fprintf(os.Stderr, "web view: open %s in your browser\n", url)
		return
	}
	exec.Command(cmd, args...).Start() //nolint:errcheck
}

func (ws *WebServer) handleMaze(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	fmt.Fprint(w, mazeSVG(ws.maze))
}

// mazeSVG generates an SVG string for the maze.
func mazeSVG(m *mazegen.Maze) string {
	const cell = 20
	const wall = 2
	const pad = wall

	svgW := m.Width*cell + 2*pad
	svgH := m.Height*cell + 2*pad

	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`, svgW, svgH)
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="white"/>`, svgW, svgH)

	strokeColor := "#1e1e1e"
	strokeW := fmt.Sprintf("%d", wall)

	line := func(x1, y1, x2, y2 int) {
		fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="%s"/>`,
			x1+pad, y1+pad, x2+pad, y2+pad, strokeColor, strokeW)
	}

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			px, py := x*cell, y*cell
			p := mazegen.Point{X: x, Y: y}
			if m.HasWall(p, mazegen.North) {
				line(px, py, px+cell, py)
			}
			if m.HasWall(p, mazegen.West) {
				line(px, py, px, py+cell)
			}
			if y == m.Height-1 && m.HasWall(p, mazegen.South) {
				line(px, py+cell, px+cell, py+cell)
			}
			if x == m.Width-1 && m.HasWall(p, mazegen.East) {
				line(px+cell, py, px+cell, py+cell)
			}
		}
	}

	// Entry marker (green circle).
	ex := m.Entry.X*cell + cell/2 + pad
	ey := m.Entry.Y*cell + cell/2 + pad
	fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="%d" fill="#3cb371"/>`, ex, ey, cell/4)

	// Exit marker (blue circle).
	xx := m.Exit.X*cell + cell/2 + pad
	xy := m.Exit.Y*cell + cell/2 + pad
	fmt.Fprintf(&sb, `<circle cx="%d" cy="%d" r="%d" fill="#4169e1"/>`, xx, xy, cell/4)

	// Solution path.
	if m.Solution != nil {
		cur := m.Entry
		for _, dir := range m.Solution {
			d := mazegen.DirDelta[dir]
			next := mazegen.Point{X: cur.X + d.X, Y: cur.Y + d.Y}
			x1 := cur.X*cell + cell/2 + pad
			y1 := cur.Y*cell + cell/2 + pad
			x2 := next.X*cell + cell/2 + pad
			y2 := next.Y*cell + cell/2 + pad
			fmt.Fprintf(&sb, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#dc143c" stroke-width="3"/>`,
				x1, y1, x2, y2)
			cur = next
		}
	}

	sb.WriteString("</svg>")
	return sb.String()
}
