package ui

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/config"
	"github.com/pr0100pr0111/AppleSoda/internal/player"
)

func TestNavigationSequentialNextPrev(t *testing.T) {
	fp := player.NewFakePlayer()
	cfg := config.Default()
	m := NewModel(fp, cfg, nil)

	testTracks := []player.Track{
		{PersistentID: "id-1", Name: "Track 1", Artist: "Artist 1", Duration: 120},
		{PersistentID: "id-2", Name: "Track 2", Artist: "Artist 2", Duration: 180},
		{PersistentID: "id-3", Name: "Track 3", Artist: "Artist 3", Duration: 200},
	}
	m.tracks = testTracks
	m.currentPLID = "pl-test"
	m.active = panelTracks
	m.selectedTrack = 0

	mod, cmd := m.handleEnter()
	m = mod.(Model)
	if cmd == nil {
		t.Fatal("expected play cmd on Enter")
	}
	if m.playingTrackID != "id-1" {
		t.Fatalf("expected playingTrackID 'id-1', got %q", m.playingTrackID)
	}

	mod, cmd = m.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	m = mod.(Model)
	if cmd == nil {
		t.Fatal("expected play cmd on Next")
	}
	if m.playingTrackID != "id-2" {
		t.Fatalf("expected playingTrackID 'id-2', got %q", m.playingTrackID)
	}

	mod, cmd = m.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	m = mod.(Model)
	if cmd == nil {
		t.Fatal("expected play cmd on Next")
	}
	if m.playingTrackID != "id-3" {
		t.Fatalf("expected playingTrackID 'id-3', got %q", m.playingTrackID)
	}

	m.interpPos = 15.0
	mod, cmd = m.Update(tea.KeyPressMsg{Code: 'p', Text: "p"})
	m = mod.(Model)
	if cmd == nil {
		t.Fatal("expected seek cmd on Prev when pos > 3")
	}
	if m.playingTrackID != "id-3" {
		t.Fatalf("expected track to remain 'id-3', got %q", m.playingTrackID)
	}

	m.interpPos = 1.0
	mod, cmd = m.Update(tea.KeyPressMsg{Code: 'p', Text: "p"})
	m = mod.(Model)
	if cmd == nil {
		t.Fatal("expected play cmd on Prev")
	}
	if m.playingTrackID != "id-2" {
		t.Fatalf("expected playingTrackID 'id-2', got %q", m.playingTrackID)
	}

	m.interpPos = 0.5
	mod, cmd = m.Update(tea.KeyPressMsg{Code: 'p', Text: "p"})
	m = mod.(Model)
	if m.playingTrackID != "id-1" {
		t.Fatalf("expected playingTrackID 'id-1', got %q", m.playingTrackID)
	}
}

func TestNavigationRepeatAllWrap(t *testing.T) {
	fp := player.NewFakePlayer()
	cfg := config.Default()
	m := NewModel(fp, cfg, nil)

	testTracks := []player.Track{
		{PersistentID: "id-1", Name: "Track 1", Duration: 100},
		{PersistentID: "id-2", Name: "Track 2", Duration: 100},
	}
	m.tracks = testTracks
	m.currentPLID = "pl-test"
	m.status.Repeat = player.RepeatAll

	m.selectedTrack = 1
	m.active = panelTracks
	mod, _ := m.handleEnter()
	m = mod.(Model)

	mod, _ = m.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	m = mod.(Model)
	if m.playingTrackID != "id-1" {
		t.Fatalf("expected wrap around to 'id-1', got %q", m.playingTrackID)
	}

	m.interpPos = 0.0
	mod, _ = m.Update(tea.KeyPressMsg{Code: 'p', Text: "p"})
	m = mod.(Model)
	if m.playingTrackID != "id-2" {
		t.Fatalf("expected wrap around to 'id-2', got %q", m.playingTrackID)
	}
}

func TestNavigationEndOfLibraryRepeatOff(t *testing.T) {
	fp := player.NewFakePlayer()
	cfg := config.Default()
	m := NewModel(fp, cfg, nil)

	testTracks := []player.Track{
		{PersistentID: "id-1", Name: "Track 1"},
		{PersistentID: "id-2", Name: "Track 2"},
	}
	m.tracks = testTracks
	m.currentPLID = "pl-test"
	m.status.Repeat = player.RepeatOff

	m.selectedTrack = 1
	m.active = panelTracks
	mod, _ := m.handleEnter()
	m = mod.(Model)

	mod, _ = m.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	m = mod.(Model)
	if !m.playlistFinished {
		t.Fatal("expected playlistFinished true when last track ends with RepeatOff")
	}
}

func TestNavigationShuffleMode(t *testing.T) {
	fp := player.NewFakePlayer()
	cfg := config.Default()
	m := NewModel(fp, cfg, nil)

	testTracks := []player.Track{
		{PersistentID: "id-1", Name: "Track 1"},
		{PersistentID: "id-2", Name: "Track 2"},
		{PersistentID: "id-3", Name: "Track 3"},
		{PersistentID: "id-4", Name: "Track 4"},
	}
	m.tracks = testTracks
	m.currentPLID = "pl-test"
	m.active = panelTracks
	m.status.Shuffle = true

	m.selectedTrack = 0
	mod, _ := m.handleEnter()
	m = mod.(Model)

	if len(m.shufflePool) != 3 {
		t.Fatalf("expected shufflePool length 3, got %d", len(m.shufflePool))
	}

	played := map[string]bool{"id-1": true}

	for i := 0; i < 3; i++ {
		mod, _ = m.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
		m = mod.(Model)
		if played[m.playingTrackID] {
			t.Fatalf("track %s was played twice before pool exhausted", m.playingTrackID)
		}
		played[m.playingTrackID] = true
	}

	if len(played) != 4 {
		t.Fatalf("expected all 4 tracks to be played in shuffle, got %d", len(played))
	}

	lastPlayed := m.playingTrackID
	m.interpPos = 0.0
	mod, _ = m.Update(tea.KeyPressMsg{Code: 'p', Text: "p"})
	m = mod.(Model)
	if m.playingTrackID == lastPlayed {
		t.Fatal("expected previous track in shuffle history, got same track")
	}
}

func TestNavigationAutoAdvance(t *testing.T) {
	fp := player.NewFakePlayer()
	cfg := config.Default()
	m := NewModel(fp, cfg, nil)

	testTracks := []player.Track{
		{PersistentID: "id-1", Name: "Track 1", Duration: 100},
		{PersistentID: "id-2", Name: "Track 2", Duration: 150},
	}
	m.tracks = testTracks
	m.currentPLID = "pl-test"
	m.active = panelTracks

	m.selectedTrack = 0
	mod, _ := m.handleEnter()
	m = mod.(Model)
	m.lastPlayActionAt = time.Now().Add(-2 * time.Second)

	stoppedStatus := player.Status{
		State:    player.StateStopped,
		HasTrack: false,
	}
	mod, cmd := m.Update(statusMsg{status: stoppedStatus})
	m = mod.(Model)

	if cmd == nil {
		t.Fatal("expected auto-advance play cmd when track stopped")
	}
	if m.playingTrackID != "id-2" {
		t.Fatalf("expected auto-advance to track 'id-2', got %q", m.playingTrackID)
	}

	m.lastPlayActionAt = time.Now().Add(-2 * time.Second)
	externalStatus := player.Status{
		State:    player.StatePlaying,
		HasTrack: true,
		Track: player.Track{
			PersistentID: "scripting-random-song",
			Name:         "Random Autoplay Song",
		},
	}
	mod, _ = m.Update(statusMsg{status: externalStatus})
	m = mod.(Model)

	if !m.playlistFinished {
		t.Fatal("expected playlistFinished true after all tracks played")
	}
}
