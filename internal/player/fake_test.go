package player

import (
	"context"
	"testing"
)

func TestFakePlayerStatus(t *testing.T) {
	ctx := context.Background()
	p := NewFakePlayer()

	st, err := p.Status(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st.State != StateStopped {
		t.Fatalf("expected StateStopped, got %v", st.State)
	}
	if st.Volume != 75 {
		t.Fatalf("expected volume 75, got %d", st.Volume)
	}
}

func TestFakePlayerControls(t *testing.T) {
	ctx := context.Background()
	p := NewFakePlayer()

	if err := p.Play(ctx); err != nil {
		t.Fatalf("play error: %v", err)
	}
	st, _ := p.Status(ctx)
	if st.State != StatePlaying {
		t.Fatalf("expected playing, got %v", st.State)
	}

	if err := p.Pause(ctx); err != nil {
		t.Fatalf("pause error: %v", err)
	}
	st, _ = p.Status(ctx)
	if st.State != StatePaused {
		t.Fatalf("expected paused, got %v", st.State)
	}

	if err := p.TogglePlayPause(ctx); err != nil {
		t.Fatalf("toggle error: %v", err)
	}
	st, _ = p.Status(ctx)
	if st.State != StatePlaying {
		t.Fatalf("expected playing, got %v", st.State)
	}

	if err := p.Next(ctx); err != nil {
		t.Fatalf("next error: %v", err)
	}

	if err := p.Previous(ctx); err != nil {
		t.Fatalf("prev error: %v", err)
	}

	if err := p.Seek(ctx, 42.0); err != nil {
		t.Fatalf("seek error: %v", err)
	}
	st, _ = p.Status(ctx)
	if st.Position < 42.0 {
		t.Fatalf("expected position >= 42.0, got %f", st.Position)
	}

	if err := p.SetVolume(ctx, 80); err != nil {
		t.Fatalf("set volume error: %v", err)
	}
	st, _ = p.Status(ctx)
	if st.Volume != 80 {
		t.Fatalf("expected volume 80, got %d", st.Volume)
	}

	if err := p.SetShuffle(ctx, true); err != nil {
		t.Fatalf("set shuffle error: %v", err)
	}
	st, _ = p.Status(ctx)
	if !st.Shuffle {
		t.Fatal("expected shuffle true")
	}

	if err := p.SetRepeat(ctx, RepeatAll); err != nil {
		t.Fatalf("set repeat error: %v", err)
	}
	st, _ = p.Status(ctx)
	if st.Repeat != RepeatAll {
		t.Fatalf("expected repeat all, got %v", st.Repeat)
	}
}

func TestFakePlayerPlaylistsAndTracks(t *testing.T) {
	ctx := context.Background()
	p := NewFakePlayer()

	pls, err := p.ListPlaylists(ctx)
	if err != nil {
		t.Fatalf("list playlists error: %v", err)
	}
	if len(pls) == 0 {
		t.Fatal("expected non-empty playlists")
	}

	tracks, total, err := p.ListTracks(ctx, pls[0].PersistentID, 0, 10)
	if err != nil {
		t.Fatalf("list tracks error: %v", err)
	}
	if total == 0 || len(tracks) == 0 {
		t.Fatal("expected tracks")
	}

	if err := p.PlayTrack(ctx, pls[0].PersistentID, tracks[0].PersistentID); err != nil {
		t.Fatalf("play track error: %v", err)
	}
	st, _ := p.Status(ctx)
	if st.Track.PersistentID != tracks[0].PersistentID {
		t.Fatalf("expected track %s, got %s", tracks[0].PersistentID, st.Track.PersistentID)
	}
}
