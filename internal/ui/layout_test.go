package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/config"
	"github.com/pr0100pr0111/AppleSoda/internal/lyrics"
	"github.com/pr0100pr0111/AppleSoda/internal/player"
)

func TestNarrowWindowLyricsVisible(t *testing.T) {
	fp := player.NewFakePlayer()
	cfg := config.Config{}
	cfg.UI.Cover = true
	cfg.UI.Lyrics = true
	cfg.Theme = "apple"
	m := NewModel(fp, cfg, nil)
	m.status = player.Status{
		HasTrack: true,
		Track: player.Track{
			PersistentID: "t01",
			Name:         "Midnight City",
			Artist:       "M83",
			Album:        "Hurry Up",
			Duration:     240,
		},
		State:    player.StatePlaying,
		Position: 15,
	}
	m.showLyrics = true
	m.lyricsResult = lyrics.Result{
		Synced: true,
		Lines: []lyrics.Line{
			{TimeSec: 5, Text: "Waiting in a car"},
			{TimeSec: 12, Text: "Waiting for a ride in the dark"},
			{TimeSec: 20, Text: "The night city grows"},
		},
	}

	testSizes := [][2]int{
		{60, 16},
		{65, 18},
		{70, 20},
		{80, 24},
		{100, 30},
		{140, 40},
	}

	for _, sz := range testSizes {
		w, h := sz[0], sz[1]
		mod, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
		curr := mod.(Model)
		v := curr.View()

		if curr.tooSmall {
			t.Fatalf("unexpected tooSmall for %dx%d", w, h)
		}

		if !strings.Contains(v.Content, "Midnight City") {
			t.Fatalf("expected track title in view for %dx%d", w, h)
		}

		if !strings.Contains(v.Content, "Waiting for a ride") {
			t.Fatalf("expected active lyrics line in view for %dx%d", w, h)
		}
	}
}

func TestNowPlayingCompactWithoutCover(t *testing.T) {
	fp := player.NewFakePlayer()
	cfg := config.Config{}
	cfg.UI.Cover = false
	cfg.UI.Lyrics = true
	cfg.Theme = "apple"
	m := NewModel(fp, cfg, nil)
	m.status = player.Status{
		HasTrack: true,
		Track: player.Track{
			PersistentID: "t02",
			Name:         "Resonance",
			Artist:       "HOME",
			Duration:     212,
		},
		State:    player.StatePlaying,
		Position: 10,
	}
	m.showLyrics = true
	m.lyricsResult = lyrics.Result{
		Plain: "Line 1 of plain lyrics\nLine 2 of plain lyrics\nLine 3 of plain lyrics",
	}

	panel := m.renderNowPlayingPanel(30, 10)
	if !strings.Contains(panel, "Resonance") {
		t.Fatalf("expected track name in panel: %s", panel)
	}
	if !strings.Contains(panel, "Line 1 of plain lyrics") {
		t.Fatalf("expected plain lyrics in panel: %s", panel)
	}
}
