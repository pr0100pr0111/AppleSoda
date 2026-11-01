package ui

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/config"
	"github.com/pr0100pr0111/AppleSoda/internal/player"
)

func TestModelInit(t *testing.T) {
	p := player.NewFakePlayer()
	cfg := config.Default()
	m := NewModel(p, cfg)

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("expected non-nil Init cmd")
	}
}

func TestModelWindowResizing(t *testing.T) {
	p := player.NewFakePlayer()
	cfg := config.Default()
	m := NewModel(p, cfg)

	sizes := [][2]int{
		{40, 12},
		{80, 24},
		{120, 40},
		{200, 60},
	}

	for _, sz := range sizes {
		w, h := sz[0], sz[1]
		mod, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
		curr := mod.(Model)
		if curr.width != w || curr.height != h {
			t.Fatalf("expected size %dx%d, got %dx%d", w, h, curr.width, curr.height)
		}
		view := curr.View()
		if w < minWidth || h < minHeight {
			if !curr.tooSmall {
				t.Fatalf("expected tooSmall for %dx%d", w, h)
			}
			if !strings.Contains(view.Content, "too small") {
				t.Fatalf("expected 'too small' in view content for %dx%d", w, h)
			}
		} else {
			if curr.tooSmall {
				t.Fatalf("expected not tooSmall for %dx%d", w, h)
			}
			if len(view.Content) == 0 {
				t.Fatalf("expected non-empty view for %dx%d", w, h)
			}
			if !strings.Contains(view.Content, "Vol:") {
				t.Fatalf("expected 'Vol:' in view content for %dx%d", w, h)
			}
		}
	}
}

func TestModelKeyNavigation(t *testing.T) {
	p := player.NewFakePlayer()
	cfg := config.Default()
	m := NewModel(p, cfg)
	mod, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = mod.(Model)

	pls, _ := p.ListPlaylists(t.Context())
	mod, _ = m.Update(playlistsMsg{playlists: pls})
	m = mod.(Model)

	if m.active != panelLibrary {
		t.Fatalf("expected panelLibrary, got %v", m.active)
	}

	mod, _ = m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	m = mod.(Model)
	if m.selectedPL != 1 {
		t.Fatalf("expected selectedPL 1, got %d", m.selectedPL)
	}

	mod, _ = m.Update(tea.KeyPressMsg{Code: 'k', Text: "k"})
	m = mod.(Model)
	if m.selectedPL != 0 {
		t.Fatalf("expected selectedPL 0, got %d", m.selectedPL)
	}

	mod, _ = m.Update(tea.KeyPressMsg{Code: '\t'})
	m = mod.(Model)
	if m.active != panelTracks {
		t.Fatalf("expected active panelTracks after Tab, got %v", m.active)
	}

	mod, _ = m.Update(tea.KeyPressMsg{Code: '\t'})
	m = mod.(Model)
	if m.active != panelNowPlaying {
		t.Fatalf("expected active panelNowPlaying after second Tab, got %v", m.active)
	}

	mod, _ = m.Update(tea.KeyPressMsg{Code: 'h', Text: "h"})
	m = mod.(Model)
	if m.active != panelTracks {
		t.Fatalf("expected panelTracks after h from now playing, got %v", m.active)
	}

	mod, _ = m.Update(tea.KeyPressMsg{Code: 'h', Text: "h"})
	m = mod.(Model)
	if m.active != panelLibrary {
		t.Fatalf("expected panelLibrary after second h, got %v", m.active)
	}

	mod, _ = m.Update(tea.KeyPressMsg{Code: 'l', Text: "l"})
	m = mod.(Model)
	if m.active != panelTracks {
		t.Fatalf("expected panelTracks after l from library, got %v", m.active)
	}
}

func TestModelToggles(t *testing.T) {
	p := player.NewFakePlayer()
	cfg := config.Default()
	m := NewModel(p, cfg)
	mod, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = mod.(Model)

	initialTheme := m.theme.Name
	mod, _ = m.Update(tea.KeyPressMsg{Code: 't', Text: "t"})
	m = mod.(Model)
	if m.theme.Name == initialTheme {
		t.Fatal("expected theme to change after 't'")
	}

	initialViz := m.vizMode
	mod, _ = m.Update(tea.KeyPressMsg{Code: 'v', Text: "v"})
	m = mod.(Model)
	if m.vizMode == initialViz {
		t.Fatal("expected vizMode to change after 'v'")
	}

	initialCover := m.showCover
	mod, _ = m.Update(tea.KeyPressMsg{Code: 'c', Text: "c"})
	m = mod.(Model)
	if m.showCover == initialCover {
		t.Fatal("expected showCover to toggle after 'c'")
	}

	initialLyrics := m.showLyrics
	mod, _ = m.Update(tea.KeyPressMsg{Code: 'L', Text: "L"})
	m = mod.(Model)
	if m.showLyrics == initialLyrics {
		t.Fatal("expected showLyrics to toggle after 'L'")
	}

	mod, _ = m.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	m = mod.(Model)
	if !m.showHelp {
		t.Fatal("expected showHelp true after '?'")
	}

	mod, _ = m.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	m = mod.(Model)
	if m.showHelp {
		t.Fatal("expected showHelp false after closing '?'")
	}
}

func TestModelFiltering(t *testing.T) {
	p := player.NewFakePlayer()
	cfg := config.Default()
	m := NewModel(p, cfg)
	mod, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = mod.(Model)

	tracks, total, _ := p.ListTracks(t.Context(), "lib", 0, 50)
	mod, _ = m.Update(tracksMsg{tracks: tracks, total: total, playlistID: "lib"})
	m = mod.(Model)

	mod, _ = m.Update(tea.KeyPressMsg{Code: '/', Text: "/"})
	m = mod.(Model)
	if !m.filterMode {
		t.Fatal("expected filterMode true after '/'")
	}

	mod, _ = m.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	m = mod.(Model)
	mod, _ = m.Update(tea.KeyPressMsg{Code: 'i', Text: "i"})
	m = mod.(Model)
	mod, _ = m.Update(tea.KeyPressMsg{Code: 'g', Text: "g"})
	m = mod.(Model)

	if m.filterText != "nig" {
		t.Fatalf("expected filterText 'nig', got %q", m.filterText)
	}

	filtered := m.currentDisplayTracks()
	for _, tr := range filtered {
		combined := strings.ToLower(tr.Name + " " + tr.Artist + " " + tr.Album)
		if !strings.Contains(combined, "nig") {
			t.Fatalf("track %s does not match filter", tr.Name)
		}
	}

	mod, _ = m.Update(tea.KeyPressMsg{Code: 0x1b})
	m = mod.(Model)
	if m.filterMode {
		t.Fatal("expected filterMode false after esc")
	}
}

func TestModelTickAndStatusInterpolation(t *testing.T) {
	p := player.NewFakePlayer()
	cfg := config.Default()
	m := NewModel(p, cfg)
	mod, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = mod.(Model)

	st := player.Status{
		State:    player.StatePlaying,
		Position: 10.0,
		Track: player.Track{
			Duration: 100.0,
		},
	}
	mod, _ = m.Update(statusMsg{status: st})
	m = mod.(Model)

	if m.interpPos != 10.0 {
		t.Fatalf("expected interpPos 10.0, got %f", m.interpPos)
	}

	mod, _ = m.Update(tickMsg(time.Now()))
	m = mod.(Model)
	if m.interpPos <= 10.0 {
		t.Fatalf("expected interpPos > 10.0 after tick, got %f", m.interpPos)
	}

	mod, cmd := m.Update(tea.KeyPressMsg{Code: ']', Text: "]"})
	m = mod.(Model)
	if cmd == nil {
		t.Fatal("expected seek cmd on ']'")
	}
}
