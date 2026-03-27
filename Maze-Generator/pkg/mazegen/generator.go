package mazegen

import "fmt"

// Generator generates the internal passages of a maze.
// Implementors must guarantee full connectivity (all cells reachable from Entry).
type Generator interface {
	Generate(m *Maze) error
}

// Options controls maze generation behaviour.
type Options struct {
	Perfect bool  // true = spanning tree only; false = add ~15% extra passages
	Seed    int64 // RNG seed for reproducibility
}

// NewGenerator returns the Generator for the given algorithm name.
// Valid names: "dfs", "prims".
func NewGenerator(algorithm string, opts Options) (Generator, error) {
	switch algorithm {
	case "dfs":
		return &DFSGenerator{opts: opts}, nil
	case "prims":
		return &PrimsGenerator{opts: opts}, nil
	default:
		return nil, fmt.Errorf("unknown algorithm %q: must be \"dfs\" or \"prims\"", algorithm)
	}
}
