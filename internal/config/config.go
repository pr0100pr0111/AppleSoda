package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Theme      string     `toml:"theme"`
	Visualizer Visualizer `toml:"visualizer"`
	UI         UI         `toml:"ui"`
	Music      Music      `toml:"music"`
}

type Visualizer struct {
	Source     string `toml:"source"`
	Mode       string `toml:"mode"`
	Pulse      bool   `toml:"pulse"`
	CavaSource string `toml:"cava_source"`
}

type UI struct {
	Cover  bool `toml:"cover"`
	Lyrics bool `toml:"lyrics"`
	Mouse  bool `toml:"mouse"`
}

type Music struct {
	AutoLaunch bool `toml:"autolaunch"`
}

func Default() Config {
	return Config{
		Theme: "apple",
		Visualizer: Visualizer{
			Source: "auto",
			Mode:   "bars",
			Pulse:  true,
		},
		UI: UI{
			Cover:  true,
			Lyrics: true,
			Mouse:  true,
		},
		Music: Music{
			AutoLaunch: false,
		},
	}
}

func configDir() (string, error) {
	if v := os.Getenv("APPLESODA_CONFIG_DIR"); v != "" {
		return v, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	return filepath.Join(home, ".config", "applesoda"), nil
}

func ConfigDir() (string, error) {
	return configDir()
}

func CacheDir() (string, error) {
	if v := os.Getenv("XDG_CACHE_HOME"); v != "" {
		return filepath.Join(v, "applesoda"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	return filepath.Join(home, ".cache", "applesoda"), nil
}

func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		dir, err := configDir()
		if err != nil {
			return cfg, err
		}
		path = filepath.Join(dir, "config.toml")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read config %s: %w", path, err)
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config %s: %w", path, err)
	}
	if err := validate(&cfg); err != nil {
		return cfg, fmt.Errorf("config %s: %w", path, err)
	}
	return cfg, nil
}

func validate(cfg *Config) error {
	switch cfg.Visualizer.Source {
	case "auto", "native", "cava", "simulated":
	default:
		return fmt.Errorf("invalid visualizer.source %q (want auto|native|cava|simulated)", cfg.Visualizer.Source)
	}
	switch cfg.Visualizer.Mode {
	case "bars", "mirror":
	default:
		return fmt.Errorf("invalid visualizer.mode %q (want bars|mirror)", cfg.Visualizer.Mode)
	}
	return nil
}
