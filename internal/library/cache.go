package library

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type cacheData struct {
	Playlists []Playlist         `json:"playlists"`
	Tracks    map[string][]Track `json:"tracks"`
}

func cacheFilePath(cacheDir string) string {
	return filepath.Join(cacheDir, "library.json")
}

func SaveCache(cacheDir string, playlists []Playlist, tracks map[string][]Track) error {
	if cacheDir == "" {
		return nil
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	data := cacheData{
		Playlists: playlists,
		Tracks:    tracks,
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal library cache: %w", err)
	}

	target := cacheFilePath(cacheDir)
	tmp := target + ".tmp"

	if err := os.WriteFile(tmp, raw, 0644); err != nil {
		return fmt.Errorf("write cache tmp: %w", err)
	}

	if err := os.Rename(tmp, target); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename cache file: %w", err)
	}

	return nil
}

func LoadCache(cacheDir string) ([]Playlist, map[string][]Track, error) {
	if cacheDir == "" {
		return nil, nil, nil
	}

	target := cacheFilePath(cacheDir)
	raw, err := os.ReadFile(target)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("read library cache: %w", err)
	}

	var data cacheData
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, nil, fmt.Errorf("parse library cache: %w", err)
	}

	if data.Tracks == nil {
		data.Tracks = make(map[string][]Track)
	}

	return data.Playlists, data.Tracks, nil
}
