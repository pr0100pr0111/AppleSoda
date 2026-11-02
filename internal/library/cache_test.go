package library

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadCache(t *testing.T) {
	dir := t.TempDir()

	pls := []Playlist{
		{PersistentID: "PL01", Name: "Favorites", TrackCount: 2},
	}
	trs := map[string][]Track{
		"PL01": {
			{PersistentID: "T01", Name: "Track A", Artist: "Artist A"},
			{PersistentID: "T02", Name: "Track B", Artist: "Artist B"},
		},
	}

	if err := SaveCache(dir, pls, trs); err != nil {
		t.Fatalf("save cache error: %v", err)
	}

	loadedPls, loadedTrs, err := LoadCache(dir)
	if err != nil {
		t.Fatalf("load cache error: %v", err)
	}

	if len(loadedPls) != 1 || loadedPls[0].Name != "Favorites" {
		t.Fatalf("unexpected loaded playlists: %+v", loadedPls)
	}

	cachedList, ok := loadedTrs["PL01"]
	if !ok || len(cachedList) != 2 {
		t.Fatalf("unexpected loaded tracks: %+v", loadedTrs)
	}
	if cachedList[0].Name != "Track A" {
		t.Fatalf("expected Track A, got %s", cachedList[0].Name)
	}
}

func TestLoadCacheMissing(t *testing.T) {
	pls, trs, err := LoadCache("/path/does/not/exist")
	if err != nil {
		t.Fatalf("unexpected error for missing cache: %v", err)
	}
	if pls != nil || trs != nil {
		t.Fatal("expected nil results for missing cache")
	}
}

func TestLoadCacheCorrupt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "library.json")
	if err := os.WriteFile(path, []byte("{corrupt json"), 0644); err != nil {
		t.Fatal(err)
	}

	_, _, err := LoadCache(dir)
	if err == nil {
		t.Fatal("expected error on corrupt cache")
	}
}
