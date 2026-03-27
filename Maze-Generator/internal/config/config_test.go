package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), "config.txt")
	os.WriteFile(f, []byte(content), 0644)
	return f
}

func TestConfig_ValidMinimal(t *testing.T) {
	f := writeConfig(t, "WIDTH=10\nHEIGHT=8\nENTRY=0,0\nEXIT=9,7\nOUTPUT_FILE=out.txt\nPERFECT=true\n")
	cfg, err := Load(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Width != 10 || cfg.Height != 8 {
		t.Errorf("dimensions wrong: %dx%d", cfg.Width, cfg.Height)
	}
	if cfg.Algorithm != "dfs" {
		t.Errorf("default algorithm should be dfs, got %q", cfg.Algorithm)
	}
}

func TestConfig_MissingRequired(t *testing.T) {
	// Missing OUTPUT_FILE
	f := writeConfig(t, "WIDTH=10\nHEIGHT=8\nENTRY=0,0\nEXIT=9,7\nPERFECT=true\n")
	_, err := Load(f)
	if err == nil {
		t.Fatal("expected error for missing OUTPUT_FILE")
	}
}

func TestConfig_EntryNotOnBorder(t *testing.T) {
	f := writeConfig(t, "WIDTH=10\nHEIGHT=8\nENTRY=3,3\nEXIT=9,7\nOUTPUT_FILE=o.txt\nPERFECT=true\n")
	_, err := Load(f)
	if err == nil {
		t.Fatal("expected error for non-border ENTRY")
	}
}

func TestConfig_EntryEqualsExit(t *testing.T) {
	f := writeConfig(t, "WIDTH=10\nHEIGHT=8\nENTRY=0,0\nEXIT=0,0\nOUTPUT_FILE=o.txt\nPERFECT=true\n")
	_, err := Load(f)
	if err == nil {
		t.Fatal("expected error for ENTRY == EXIT")
	}
}

func TestConfig_InvalidAlgorithm(t *testing.T) {
	f := writeConfig(t, "WIDTH=10\nHEIGHT=8\nENTRY=0,0\nEXIT=9,7\nOUTPUT_FILE=o.txt\nPERFECT=true\nALGORITHM=dijkstra\n")
	_, err := Load(f)
	if err == nil {
		t.Fatal("expected error for invalid ALGORITHM")
	}
}

func TestConfig_CommentsIgnored(t *testing.T) {
	f := writeConfig(t, "# comment\nWIDTH=10\n# another\nHEIGHT=8\nENTRY=0,0\nEXIT=9,7\nOUTPUT_FILE=o.txt\nPERFECT=false\n")
	cfg, err := Load(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Perfect {
		t.Error("expected Perfect=false")
	}
}

func TestConfig_UnknownKeyIgnored(t *testing.T) {
	f := writeConfig(t, "WIDTH=10\nHEIGHT=8\nENTRY=0,0\nEXIT=9,7\nOUTPUT_FILE=o.txt\nPERFECT=true\nFOO=bar\n")
	_, err := Load(f)
	if err != nil {
		t.Fatalf("unknown key should be ignored, got error: %v", err)
	}
}

func TestConfig_ExitNotOnBorder(t *testing.T) {
	f := writeConfig(t, "WIDTH=10\nHEIGHT=8\nENTRY=0,0\nEXIT=5,5\nOUTPUT_FILE=o.txt\nPERFECT=true\n")
	_, err := Load(f)
	if err == nil {
		t.Fatal("expected error for non-border EXIT")
	}
}
