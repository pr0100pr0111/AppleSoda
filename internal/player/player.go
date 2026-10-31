package player

import "context"

type PlayerState int

const (
	StateStopped PlayerState = iota
	StatePlaying
	StatePaused
)

func (s PlayerState) String() string {
	switch s {
	case StatePlaying:
		return "playing"
	case StatePaused:
		return "paused"
	default:
		return "stopped"
	}
}

type RepeatMode int

const (
	RepeatOff RepeatMode = iota
	RepeatAll
	RepeatOne
)

func (r RepeatMode) String() string {
	switch r {
	case RepeatAll:
		return "all"
	case RepeatOne:
		return "one"
	default:
		return "off"
	}
}

func (r RepeatMode) Next() RepeatMode {
	return (r + 1) % 3
}

type Track struct {
	PersistentID string
	Name         string
	Artist       string
	Album        string
	Duration     float64
	TrackNumber  int
	DiscNumber   int
	Genre        string
	Year         int
}

type Playlist struct {
	PersistentID string
	Name         string
	TrackCount   int
	Kind         string
}

type Status struct {
	State      PlayerState
	Position   float64
	Track      Track
	Volume     int
	Shuffle    bool
	Repeat     RepeatMode
	HasTrack   bool
	PlaylistID string
}

type Player interface {
	Status(ctx context.Context) (Status, error)
	Play(ctx context.Context) error
	Pause(ctx context.Context) error
	TogglePlayPause(ctx context.Context) error
	Next(ctx context.Context) error
	Previous(ctx context.Context) error
	Seek(ctx context.Context, position float64) error
	SetVolume(ctx context.Context, volume int) error
	SetShuffle(ctx context.Context, on bool) error
	SetRepeat(ctx context.Context, mode RepeatMode) error
	PlayTrack(ctx context.Context, playlistID string, trackID string) error
	ListPlaylists(ctx context.Context) ([]Playlist, error)
	ListTracks(ctx context.Context, playlistID string, offset, limit int) ([]Track, int, error)
	Artwork(ctx context.Context) ([]byte, string, error)
	IsRunning(ctx context.Context) bool
	Launch(ctx context.Context) error
}
