package ui

import (
	"strings"
	"testing"

	"github.com/pr0100pr0111/AppleSoda/internal/lyrics"
	"github.com/pr0100pr0111/AppleSoda/internal/player"
	"github.com/pr0100pr0111/AppleSoda/internal/theme"
)

func TestRenderLyricsNoTrack(t *testing.T) {
	m := Model{
		theme:  theme.Get("apple"),
		status: player.Status{HasTrack: false},
	}
	out := m.renderLyrics(40, 10)
	if out != "" {
		t.Fatalf("expected empty string, got %q", out)
	}
}

func TestRenderLyricsLoading(t *testing.T) {
	m := Model{
		theme:         theme.Get("apple"),
		status:        player.Status{HasTrack: true},
		lyricsLoading: true,
	}
	out := m.renderLyrics(40, 10)
	if !strings.Contains(out, "Loading lyrics...") {
		t.Fatalf("expected loading indicator, got %q", out)
	}
}

func TestRenderLyricsNotFound(t *testing.T) {
	m := Model{
		theme:         theme.Get("apple"),
		status:        player.Status{HasTrack: true},
		lyricsLoading: false,
		lyricsResult:  lyrics.Result{},
	}
	out := m.renderLyrics(40, 10)
	if !strings.Contains(out, "Текст не найден") {
		t.Fatalf("expected 'Текст не найден', got %q", out)
	}
}

func TestRenderSyncedLyrics(t *testing.T) {
	lines := []lyrics.Line{
		{TimeSec: 5, Text: "Line 1"},
		{TimeSec: 10, Text: "Line 2"},
		{TimeSec: 15, Text: "Line 3"},
	}
	m := Model{
		theme:  theme.Get("apple"),
		status: player.Status{HasTrack: true},
		lyricsResult: lyrics.Result{
			Synced: true,
			Lines:  lines,
		},
		interpPos: 11,
	}

	out := m.renderLyrics(40, 5)
	if !strings.Contains(out, "▶ Line 2") {
		t.Fatalf("expected active line Line 2, got %q", out)
	}
	if !strings.Contains(out, "Line 1") {
		t.Fatalf("expected context line Line 1, got %q", out)
	}
}

func TestRenderPlainLyrics(t *testing.T) {
	plain := "Hello world\nThis is a song\nWith plain lyrics"
	m := Model{
		theme:  theme.Get("apple"),
		status: player.Status{HasTrack: true, Track: player.Track{Duration: 100}},
		lyricsResult: lyrics.Result{
			Synced: false,
			Plain:  plain,
		},
		interpPos: 10,
	}

	out := m.renderLyrics(40, 5)
	if !strings.Contains(out, "Hello world") {
		t.Fatalf("expected plain lyrics content, got %q", out)
	}
}

func TestLyricsMsgHandling(t *testing.T) {
	m := Model{
		status: player.Status{
			HasTrack: true,
			Track:    player.Track{PersistentID: "t-100"},
		},
		lyricsLoading: true,
	}

	res := lyrics.Result{
		Plain: "Some lyrics",
	}

	newM, _ := m.Update(lyricsMsg{
		persistentID: "t-100",
		result:       res,
	})

	updated := newM.(Model)
	if updated.lyricsLoading {
		t.Fatal("expected lyricsLoading to be false")
	}
	if updated.lyricsTrackID != "t-100" {
		t.Fatalf("expected lyricsTrackID t-100, got %s", updated.lyricsTrackID)
	}
	if updated.lyricsResult.Plain != "Some lyrics" {
		t.Fatalf("expected plain lyrics to be set")
	}
}
