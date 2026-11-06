package lyrics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	cacheTTL    = 7 * 24 * time.Hour
	notFoundTTL = 24 * time.Hour
)

type CachedResult struct {
	Result
	FetchedAt time.Time
}

func cacheDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".cache", "applesoda", "lyrics")
}

func cachePath(persistentID string) string {
	return filepath.Join(cacheDir(), persistentID+".json")
}

func Load(persistentID string) (Result, bool) {
	p := cachePath(persistentID)
	data, err := os.ReadFile(p)
	if err != nil {
		return Result{}, false
	}

	var cr CachedResult
	if err := json.Unmarshal(data, &cr); err != nil {
		return Result{}, false
	}

	isEmpty := !cr.Synced && cr.Plain == "" && len(cr.Lines) == 0
	ttl := cacheTTL
	if isEmpty {
		ttl = notFoundTTL
	}

	if time.Since(cr.FetchedAt) > ttl {
		return Result{}, false
	}

	return cr.Result, true
}

func Save(persistentID string, r Result) {
	dir := cacheDir()
	if dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}

	cr := CachedResult{
		Result:    r,
		FetchedAt: time.Now(),
	}
	data, err := json.Marshal(cr)
	if err != nil {
		return
	}
	_ = os.WriteFile(cachePath(persistentID), data, 0o644)
}
