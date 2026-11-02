package library

import (
	"testing"
)

func TestStoreInMemory(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)

	pls := []Playlist{
		{PersistentID: "1", Name: "Rock"},
		{PersistentID: "2", Name: "Pop"},
	}
	s.SetPlaylists(pls)

	gotPls := s.GetPlaylists()
	if len(gotPls) != 2 {
		t.Fatalf("expected 2 playlists, got %d", len(gotPls))
	}

	tracks := []Track{
		{PersistentID: "T1", Name: "Song 1"},
	}
	s.SetTracks("1", tracks)

	gotTracks, ok := s.GetTracks("1")
	if !ok || len(gotTracks) != 1 {
		t.Fatalf("expected 1 track, got %v (%t)", gotTracks, ok)
	}

	if s.TrackCount("1") != 1 {
		t.Fatalf("expected track count 1, got %d", s.TrackCount("1"))
	}
	if s.TrackCount("nonexistent") != 0 {
		t.Fatalf("expected track count 0, got %d", s.TrackCount("nonexistent"))
	}
}

func TestStorePersistence(t *testing.T) {
	dir := t.TempDir()
	s1 := NewStore(dir)

	s1.SetPlaylists([]Playlist{{PersistentID: "A", Name: "Alpha"}})
	s1.SetTracks("A", []Track{{PersistentID: "TA1", Name: "Track Alpha 1"}})

	s2 := NewStore(dir)
	if err := s2.LoadFromDisk(); err != nil {
		t.Fatalf("load from disk error: %v", err)
	}

	pls := s2.GetPlaylists()
	if len(pls) != 1 || pls[0].Name != "Alpha" {
		t.Fatalf("unexpected persisted playlists: %+v", pls)
	}

	trs, ok := s2.GetTracks("A")
	if !ok || len(trs) != 1 || trs[0].Name != "Track Alpha 1" {
		t.Fatalf("unexpected persisted tracks: %+v", trs)
	}
}
