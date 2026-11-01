package player

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestParsePlayerState(t *testing.T) {
	if ParsePlayerState("playing") != StatePlaying {
		t.Fatal("expected StatePlaying")
	}
	if ParsePlayerState("paused") != StatePaused {
		t.Fatal("expected StatePaused")
	}
	if ParsePlayerState("stopped") != StateStopped {
		t.Fatal("expected StateStopped")
	}
	if ParsePlayerState("unknown") != StateStopped {
		t.Fatal("expected fallback StateStopped")
	}
}

func TestParseRepeatMode(t *testing.T) {
	if ParseRepeatMode("all") != RepeatAll {
		t.Fatal("expected RepeatAll")
	}
	if ParseRepeatMode("one") != RepeatOne {
		t.Fatal("expected RepeatOne")
	}
	if ParseRepeatMode("off") != RepeatOff {
		t.Fatal("expected RepeatOff")
	}
}

func TestJXAStatusJSONParsing(t *testing.T) {
	fixture := []byte(`{
		"state": "playing",
		"position": 42.5,
		"volume": 85,
		"shuffle": true,
		"repeat": "all",
		"track": {
			"name": "Sunset Drive",
			"artist": "The Midnight",
			"album": "Endless Summer",
			"duration": 210.0,
			"persistentID": "ABCD1234EF567890"
		},
		"playlist": {
			"name": "Favorites",
			"persistentID": "PL001"
		}
	}`)

	var res jxaStatusResponse
	if err := json.Unmarshal(fixture, &res); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	st := Status{
		State:        ParsePlayerState(res.State),
		Position:     res.Position,
		Volume:       res.Volume,
		Shuffle:      res.Shuffle,
		Repeat:       ParseRepeatMode(res.Repeat),
		PlaylistID:   res.Playlist.PersistentID,
		PlaylistName: res.Playlist.Name,
	}

	if res.Track != nil {
		st.HasTrack = true
		st.Track = Track{
			PersistentID: res.Track.PersistentID,
			Name:         res.Track.Name,
			Artist:       res.Track.Artist,
			Album:        res.Track.Album,
			Duration:     res.Track.Duration,
		}
	}

	if st.State != StatePlaying {
		t.Fatalf("expected StatePlaying, got %v", st.State)
	}
	if st.Position != 42.5 {
		t.Fatalf("expected position 42.5, got %f", st.Position)
	}
	if st.Volume != 85 {
		t.Fatalf("expected volume 85, got %d", st.Volume)
	}
	if !st.Shuffle {
		t.Fatal("expected shuffle true")
	}
	if st.Repeat != RepeatAll {
		t.Fatalf("expected repeat all, got %v", st.Repeat)
	}
	if !st.HasTrack || st.Track.Name != "Sunset Drive" {
		t.Fatalf("expected track Sunset Drive, got %s", st.Track.Name)
	}
	if st.PlaylistName != "Favorites" {
		t.Fatalf("expected playlist Favorites, got %s", st.PlaylistName)
	}
}

func TestJXAPlaylistsJSONParsing(t *testing.T) {
	fixture := []byte(`[
		{"name": "Library", "id": "0001"},
		{"name": "Favorites", "id": "0002"}
	]`)

	var rawList []jxaPlaylist
	if err := json.Unmarshal(fixture, &rawList); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(rawList) != 2 {
		t.Fatalf("expected 2 playlists, got %d", len(rawList))
	}
	if rawList[0].Name != "Library" || rawList[0].ID != "0001" {
		t.Fatalf("unexpected playlist item: %+v", rawList[0])
	}
}

func TestJXATracksJSONParsing(t *testing.T) {
	fixture := []byte(`{
		"total": 100,
		"tracks": [
			{
				"name": "Track 1",
				"artist": "Artist 1",
				"album": "Album 1",
				"duration": 180.0,
				"persistentID": "T001",
				"index": 1
			}
		]
	}`)

	var res jxaTracksResponse
	if err := json.Unmarshal(fixture, &res); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if res.Total != 100 || len(res.Tracks) != 1 {
		t.Fatalf("expected total 100 and 1 track, got %d and %d", res.Total, len(res.Tracks))
	}
	if res.Tracks[0].Name != "Track 1" || res.Tracks[0].PersistentID != "T001" {
		t.Fatalf("unexpected track: %+v", res.Tracks[0])
	}
}

func TestMusicAppLiveIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := NewMusicAppPlayer(false)
	if !p.IsRunning(ctx) {
		t.Skip("Music.app is not running, skipping live integration test")
	}

	st, err := p.Status(ctx)
	if err != nil {
		t.Skipf("cannot query live status (permissions or busy): %v", err)
	}
	if st.Volume < 0 || st.Volume > 100 {
		t.Fatalf("unexpected live volume: %d", st.Volume)
	}

	pls, err := p.ListPlaylists(ctx)
	if err != nil {
		t.Skipf("cannot list playlists: %v", err)
	}
	if len(pls) == 0 {
		t.Fatal("expected at least 1 playlist in live Music.app")
	}
}
