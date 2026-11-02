package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Theme != "apple" {
		t.Fatalf("want theme apple, got %s", cfg.Theme)
	}
	if cfg.Visualizer.Source != "auto" {
		t.Fatalf("want source auto, got %s", cfg.Visualizer.Source)
	}
}

func TestLoadMissing(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.toml")
	if err != nil {
		t.Fatalf("unexpected error for missing config: %v", err)
	}
	if cfg.Theme != "apple" {
		t.Fatalf("want default theme, got %s", cfg.Theme)
	}
}

func TestLoadValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	data := []byte(`theme = "nord"
[visualizer]
source = "simulated"
mode = "mirror"
pulse = false
`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Theme != "nord" {
		t.Fatalf("want nord, got %s", cfg.Theme)
	}
	if cfg.Visualizer.Source != "simulated" {
		t.Fatalf("want simulated, got %s", cfg.Visualizer.Source)
	}
	if cfg.Visualizer.Mode != "mirror" {
		t.Fatalf("want mirror, got %s", cfg.Visualizer.Mode)
	}
}

func TestLoadInvalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	data := []byte(`[visualizer]
source = "invalid"
`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid source")
	}
}

func TestCacheDir(t *testing.T) {
	dir, err := CacheDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir == "" {
		t.Fatal("expected non-empty cache dir")
	}
}
