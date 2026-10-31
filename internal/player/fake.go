package player

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

type FakePlayer struct {
	mu         sync.Mutex
	state      PlayerState
	position   float64
	volume     int
	shuffle    bool
	repeat     RepeatMode
	trackIdx   int
	playlistID string
	playlists  []Playlist
	tracks     map[string][]Track
	startedAt  time.Time
}

func NewFakePlayer() *FakePlayer {
	fp := &FakePlayer{
		state:      StateStopped,
		volume:     75,
		playlistID: "lib",
	}
	fp.playlists = []Playlist{
		{PersistentID: "lib", Name: "Library", TrackCount: 0, Kind: "library"},
		{PersistentID: "fav", Name: "Favorites", TrackCount: 0, Kind: "user"},
		{PersistentID: "rock", Name: "Rock Classics", TrackCount: 0, Kind: "user"},
		{PersistentID: "chill", Name: "Chill Vibes", TrackCount: 0, Kind: "user"},
		{PersistentID: "recent", Name: "Recently Added", TrackCount: 0, Kind: "smart"},
		{PersistentID: "top", Name: "Top 25 Most Played", TrackCount: 0, Kind: "smart"},
	}

	artists := []string{
		"The Midnight", "FM-84", "Gunship", "Timecop1983", "Kavinsky",
		"Perturbator", "Carpenter Brut", "Dance With The Dead", "Lazerhawk", "Mitch Murder",
		"Com Truise", "Tycho", "HOME", "Electric Youth", "College",
	}
	albums := []string{
		"Endless Summer", "Atlas", "Dark All Day", "Night Driving", "OutRun",
		"Dangerous Days", "Leather Teeth", "Loved to Death", "Redline", "Current Events",
		"Iteration", "Dive", "Odyssey", "Innerworld", "Northern Council",
	}

	var libTracks []Track
	for i := 0; i < 200; i++ {
		t := Track{
			PersistentID: fmt.Sprintf("t%04d", i),
			Name:         fmt.Sprintf("Track %d - %s", i+1, demoTrackName(i)),
			Artist:       artists[i%len(artists)],
			Album:        albums[i%len(albums)],
			Duration:     180 + float64(i%120)*2,
			TrackNumber:  (i % 12) + 1,
			DiscNumber:   1,
			Year:         2018 + (i % 7),
		}
		libTracks = append(libTracks, t)
	}

	fp.tracks = map[string][]Track{
		"lib":    libTracks,
		"fav":    libTracks[:25],
		"rock":   libTracks[30:60],
		"chill":  libTracks[60:85],
		"recent": libTracks[150:200],
		"top":    libTracks[:25],
	}
	for i := range fp.playlists {
		if ts, ok := fp.tracks[fp.playlists[i].PersistentID]; ok {
			fp.playlists[i].TrackCount = len(ts)
		}
	}
	return fp
}

func demoTrackName(i int) string {
	names := []string{
		"Neon Lights", "Sunset Drive", "Midnight City", "Digital Love", "Starlight",
		"Turbo Killer", "Nightcall", "Acid Rain", "Overdrive", "Electric Dreams",
		"Le Perv", "Lost Years", "Resonance", "Shadows", "A Real Hero",
		"Tech Noir", "Roller Mobster", "Skull and Shark", "King of the Streets", "Dead of Night",
	}
	return names[i%len(names)]
}

func (fp *FakePlayer) currentTrack() Track {
	tracks := fp.tracks[fp.playlistID]
	if len(tracks) == 0 {
		return Track{}
	}
	idx := fp.trackIdx % len(tracks)
	return tracks[idx]
}

func (fp *FakePlayer) Status(_ context.Context) (Status, error) {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	pos := fp.position
	if fp.state == StatePlaying && !fp.startedAt.IsZero() {
		pos += time.Since(fp.startedAt).Seconds()
	}
	t := fp.currentTrack()
	if t.Duration > 0 {
		pos = math.Mod(pos, t.Duration)
	}

	return Status{
		State:      fp.state,
		Position:   pos,
		Track:      t,
		Volume:     fp.volume,
		Shuffle:    fp.shuffle,
		Repeat:     fp.repeat,
		HasTrack:   t.PersistentID != "",
		PlaylistID: fp.playlistID,
	}, nil
}

func (fp *FakePlayer) Play(_ context.Context) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	fp.state = StatePlaying
	fp.startedAt = time.Now()
	return nil
}

func (fp *FakePlayer) Pause(_ context.Context) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	if fp.state == StatePlaying && !fp.startedAt.IsZero() {
		fp.position += time.Since(fp.startedAt).Seconds()
	}
	fp.state = StatePaused
	fp.startedAt = time.Time{}
	return nil
}

func (fp *FakePlayer) TogglePlayPause(ctx context.Context) error {
	fp.mu.Lock()
	s := fp.state
	fp.mu.Unlock()
	if s == StatePlaying {
		return fp.Pause(ctx)
	}
	return fp.Play(ctx)
}

func (fp *FakePlayer) Next(_ context.Context) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	tracks := fp.tracks[fp.playlistID]
	if len(tracks) > 0 {
		fp.trackIdx = (fp.trackIdx + 1) % len(tracks)
	}
	fp.position = 0
	fp.startedAt = time.Now()
	return nil
}

func (fp *FakePlayer) Previous(_ context.Context) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	tracks := fp.tracks[fp.playlistID]
	if len(tracks) > 0 {
		fp.trackIdx = (fp.trackIdx - 1 + len(tracks)) % len(tracks)
	}
	fp.position = 0
	fp.startedAt = time.Now()
	return nil
}

func (fp *FakePlayer) Seek(_ context.Context, position float64) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	fp.position = position
	if fp.state == StatePlaying {
		fp.startedAt = time.Now()
	}
	return nil
}

func (fp *FakePlayer) SetVolume(_ context.Context, volume int) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	if volume < 0 {
		volume = 0
	}
	if volume > 100 {
		volume = 100
	}
	fp.volume = volume
	return nil
}

func (fp *FakePlayer) SetShuffle(_ context.Context, on bool) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	fp.shuffle = on
	return nil
}

func (fp *FakePlayer) SetRepeat(_ context.Context, mode RepeatMode) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	fp.repeat = mode
	return nil
}

func (fp *FakePlayer) PlayTrack(_ context.Context, playlistID string, trackID string) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	if playlistID != "" {
		fp.playlistID = playlistID
	}
	tracks := fp.tracks[fp.playlistID]
	for i, t := range tracks {
		if t.PersistentID == trackID {
			fp.trackIdx = i
			break
		}
	}
	fp.position = 0
	fp.state = StatePlaying
	fp.startedAt = time.Now()
	return nil
}

func (fp *FakePlayer) ListPlaylists(_ context.Context) ([]Playlist, error) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	result := make([]Playlist, len(fp.playlists))
	copy(result, fp.playlists)
	return result, nil
}

func (fp *FakePlayer) ListTracks(_ context.Context, playlistID string, offset, limit int) ([]Track, int, error) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	pid := playlistID
	if pid == "" {
		pid = "lib"
	}
	tracks := fp.tracks[pid]
	total := len(tracks)
	if offset >= total {
		return nil, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	result := make([]Track, end-offset)
	copy(result, tracks[offset:end])
	return result, total, nil
}

func (fp *FakePlayer) Artwork(_ context.Context) ([]byte, string, error) {
	return nil, "", nil
}

func (fp *FakePlayer) IsRunning(_ context.Context) bool {
	return true
}

func (fp *FakePlayer) Launch(_ context.Context) error {
	return nil
}
