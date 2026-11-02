package library

import (
	"testing"
)

func TestFilterTracks(t *testing.T) {
	tracks := []Track{
		{Name: "Sunset Drive", Artist: "The Midnight", Album: "Endless Summer"},
		{Name: "Neon Lights", Artist: "FM-84", Album: "Atlas"},
		{Name: "Платье в точку", Artist: "АлоэВера", Album: "Сделаем вид"},
		{Name: "Дискац", Artist: "АлоэВера", Album: "Сделаем вид"},
		{Name: "Resonance", Artist: "HOME", Album: "Odyssey"},
	}

	all := FilterTracks(tracks, "")
	if len(all) != len(tracks) {
		t.Fatalf("expected all %d tracks, got %d", len(tracks), len(all))
	}

	byArtist := FilterTracks(tracks, "midnight")
	if len(byArtist) != 1 || byArtist[0].Name != "Sunset Drive" {
		t.Fatalf("expected Sunset Drive, got %+v", byArtist)
	}

	byCyrillic := FilterTracks(tracks, "алоэвера")
	if len(byCyrillic) != 2 {
		t.Fatalf("expected 2 tracks by АлоэВера, got %d", len(byCyrillic))
	}

	multiTerm := FilterTracks(tracks, "платье вид")
	if len(multiTerm) != 1 || multiTerm[0].Name != "Платье в точку" {
		t.Fatalf("expected Платье в точку, got %+v", multiTerm)
	}

	noMatch := FilterTracks(tracks, "nonexistent")
	if len(noMatch) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(noMatch))
	}
}

func TestFilterIndices(t *testing.T) {
	tracks := []Track{
		{Name: "One", Artist: "A", Album: "X"},
		{Name: "Two", Artist: "B", Album: "Y"},
		{Name: "Three", Artist: "A", Album: "Z"},
	}

	if FilterIndices(tracks, "") != nil {
		t.Fatal("expected nil indices for empty query")
	}

	idx := FilterIndices(tracks, "a")
	if len(idx) != 2 || idx[0] != 0 || idx[1] != 2 {
		t.Fatalf("expected [0 2], got %v", idx)
	}
}
