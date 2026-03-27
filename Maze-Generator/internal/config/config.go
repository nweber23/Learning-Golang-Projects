package config

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/niklaswebde/maze-generator/pkg/mazegen"
)

// Config holds all parsed and validated maze configuration.
type Config struct {
	Width, Height int
	Entry, Exit   mazegen.Point
	OutputFile    string
	Perfect       bool
	Seed          int64
	Algorithm     string // "dfs" or "prims"
	WebView       bool
	OutputImage   string // empty = no PNG output
}

// Load parses a KEY=VALUE config file and validates all fields.
// Returns a descriptive error on any validation failure.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open config file %q: %w", path, err)
	}
	defer f.Close()

	kv := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Strip inline comments.
		if idx := strings.Index(line, " #"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // ignore malformed lines
		}
		kv[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	cfg := &Config{
		Algorithm: "dfs",       // default
		Seed:      rand.Int63(), // random if SEED not specified
	}

	// Required keys.
	required := []string{"WIDTH", "HEIGHT", "ENTRY", "EXIT", "OUTPUT_FILE", "PERFECT"}
	for _, key := range required {
		if _, ok := kv[key]; !ok {
			return nil, fmt.Errorf("config: missing required key %q", key)
		}
	}

	// WIDTH
	w, err := strconv.Atoi(kv["WIDTH"])
	if err != nil || w < 3 {
		return nil, fmt.Errorf("config: WIDTH must be an integer >= 3, got %q", kv["WIDTH"])
	}
	cfg.Width = w

	// HEIGHT
	h, err := strconv.Atoi(kv["HEIGHT"])
	if err != nil || h < 3 {
		return nil, fmt.Errorf("config: HEIGHT must be an integer >= 3, got %q", kv["HEIGHT"])
	}
	cfg.Height = h

	// ENTRY
	entry, err := parsePoint(kv["ENTRY"])
	if err != nil {
		return nil, fmt.Errorf("config: invalid ENTRY %q: %w", kv["ENTRY"], err)
	}
	if !inBounds(entry, w, h) {
		return nil, fmt.Errorf("config: ENTRY (%d,%d) is out of bounds for %dx%d maze", entry.X, entry.Y, w, h)
	}
	if !onBorder(entry, w, h) {
		return nil, fmt.Errorf("config: ENTRY (%d,%d) must be on the outer border of the maze", entry.X, entry.Y)
	}
	cfg.Entry = entry

	// EXIT
	exit, err := parsePoint(kv["EXIT"])
	if err != nil {
		return nil, fmt.Errorf("config: invalid EXIT %q: %w", kv["EXIT"], err)
	}
	if !inBounds(exit, w, h) {
		return nil, fmt.Errorf("config: EXIT (%d,%d) is out of bounds for %dx%d maze", exit.X, exit.Y, w, h)
	}
	if !onBorder(exit, w, h) {
		return nil, fmt.Errorf("config: EXIT (%d,%d) must be on the outer border of the maze", exit.X, exit.Y)
	}
	if entry == exit {
		return nil, fmt.Errorf("config: ENTRY and EXIT must be different points")
	}
	cfg.Exit = exit

	// OUTPUT_FILE
	cfg.OutputFile = kv["OUTPUT_FILE"]

	// PERFECT
	switch strings.ToLower(kv["PERFECT"]) {
	case "true":
		cfg.Perfect = true
	case "false":
		cfg.Perfect = false
	default:
		return nil, fmt.Errorf("config: PERFECT must be \"true\" or \"false\", got %q", kv["PERFECT"])
	}

	// Optional: SEED
	if s, ok := kv["SEED"]; ok {
		seed, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("config: invalid SEED %q: must be an integer", s)
		}
		cfg.Seed = seed
	}

	// Optional: ALGORITHM
	if a, ok := kv["ALGORITHM"]; ok {
		switch strings.ToLower(a) {
		case "dfs", "prims":
			cfg.Algorithm = strings.ToLower(a)
		default:
			return nil, fmt.Errorf("config: ALGORITHM must be \"dfs\" or \"prims\", got %q", a)
		}
	}

	// Optional: WEB_VIEW
	if v, ok := kv["WEB_VIEW"]; ok {
		cfg.WebView = strings.ToLower(v) == "true"
	}

	// Optional: OUTPUT_IMAGE
	if img, ok := kv["OUTPUT_IMAGE"]; ok {
		cfg.OutputImage = img
	}

	return cfg, nil
}

func parsePoint(s string) (mazegen.Point, error) {
	parts := strings.SplitN(s, ",", 2)
	if len(parts) != 2 {
		return mazegen.Point{}, fmt.Errorf("expected x,y")
	}
	x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return mazegen.Point{}, fmt.Errorf("invalid x: %w", err)
	}
	y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return mazegen.Point{}, fmt.Errorf("invalid y: %w", err)
	}
	return mazegen.Point{X: x, Y: y}, nil
}

func inBounds(p mazegen.Point, w, h int) bool {
	return p.X >= 0 && p.X < w && p.Y >= 0 && p.Y < h
}

func onBorder(p mazegen.Point, w, h int) bool {
	return p.X == 0 || p.X == w-1 || p.Y == 0 || p.Y == h-1
}
