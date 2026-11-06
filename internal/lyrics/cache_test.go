package lyrics

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func unmarshalForTest(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func marshalForTest(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	origCacheDir := os.Getenv("HOME")
	_ = origCacheDir

	origHome, _ := os.UserHomeDir()
	t.Setenv("HOME", dir)
	defer t.Setenv("HOME", origHome)

	r := Result{
		Lines:  []Line{{TimeSec: 1.0, Text: "hello"}, {TimeSec: 2.0, Text: "world"}},
		Plain:  "hello\nworld",
		Synced: true,
	}

	Save("testtrack123", r)

	loaded, ok := Load("testtrack123")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if !loaded.Synced {
		t.Error("expected Synced=true")
	}
	if loaded.Plain != "hello\nworld" {
		t.Errorf("unexpected plain: %q", loaded.Plain)
	}
	if len(loaded.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(loaded.Lines))
	}
}

func TestLoadMiss(t *testing.T) {
	dir := t.TempDir()
	origHome, _ := os.UserHomeDir()
	t.Setenv("HOME", dir)
	defer t.Setenv("HOME", origHome)

	_, ok := Load("nonexistent")
	if ok {
		t.Error("expected cache miss for nonexistent key")
	}
}

func TestLoadExpiredRegular(t *testing.T) {
	dir := t.TempDir()
	origHome, _ := os.UserHomeDir()
	t.Setenv("HOME", dir)
	defer t.Setenv("HOME", origHome)

	r := Result{Plain: "some lyrics", Synced: false}
	Save("expiredtrack", r)

	p := cachePath("expiredtrack")
	data, _ := os.ReadFile(p)

	var cr CachedResult
	_ = unmarshalForTest(data, &cr)
	cr.FetchedAt = time.Now().Add(-(8 * 24 * time.Hour))
	newData, _ := marshalForTest(cr)
	_ = os.WriteFile(p, newData, 0o644)

	_, ok := Load("expiredtrack")
	if ok {
		t.Error("expected cache miss for expired entry")
	}
}

func TestLoadExpiredNotFound(t *testing.T) {
	dir := t.TempDir()
	origHome, _ := os.UserHomeDir()
	t.Setenv("HOME", dir)
	defer t.Setenv("HOME", origHome)

	r := Result{}
	Save("notfoundtrack", r)

	p := cachePath("notfoundtrack")
	data, _ := os.ReadFile(p)

	var cr CachedResult
	_ = unmarshalForTest(data, &cr)
	cr.FetchedAt = time.Now().Add(-(2 * 24 * time.Hour))
	newData, _ := marshalForTest(cr)
	_ = os.WriteFile(p, newData, 0o644)

	_, ok := Load("notfoundtrack")
	if ok {
		t.Error("expected cache miss for expired not-found entry")
	}
}

func TestLoadNotFoundWithinTTL(t *testing.T) {
	dir := t.TempDir()
	origHome, _ := os.UserHomeDir()
	t.Setenv("HOME", dir)
	defer t.Setenv("HOME", origHome)

	r := Result{}
	Save("notfoundrecent", r)

	_, ok := Load("notfoundrecent")
	if !ok {
		t.Error("expected cache hit for recent not-found entry within 1-day TTL")
	}
}
