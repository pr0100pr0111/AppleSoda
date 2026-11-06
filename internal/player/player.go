package player

import (
	"context"
	"errors"

	"github.com/pr0100pr0111/AppleSoda/internal/library"
)

var (
	ErrMusicAppNotRunning   = errors.New("Music.app is not running (press 'o' to open)")
	ErrAutomationPermission = errors.New("Automation permission denied (System Settings → Privacy & Security → Automation → Music)")
	ErrArtworkNotFound      = errors.New("artwork not found")
	ErrTrackUnplayable      = errors.New("track cannot be played")
)

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

func ParsePlayerState(s string) PlayerState {
	switch s {
	case "playing":
		return StatePlaying
	case "paused":
		return StatePaused
	default:
		return StateStopped
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

func ParseRepeatMode(s string) RepeatMode {
	switch s {
	case "all":
		return RepeatAll
	case "one":
		return RepeatOne
	default:
		return RepeatOff
	}
}

func (r RepeatMode) Next() RepeatMode {
	return (r + 1) % 3
}

type Track = library.Track
type Playlist = library.Playlist

type Status struct {
	State        PlayerState
	Position     float64
	Track        Track
	Volume       int
	Shuffle      bool
	Repeat       RepeatMode
	HasTrack     bool
	PlaylistID   string
	PlaylistName string
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
	Lyrics(ctx context.Context) (string, error)
	IsRunning(ctx context.Context) bool
	Launch(ctx context.Context) error
}
