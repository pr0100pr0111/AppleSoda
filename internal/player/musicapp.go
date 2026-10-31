package player

import (
	"context"
	"errors"
	"os/exec"
)

var ErrMusicAppUnavailable = errors.New("music.app integration not initialized")

type MusicAppPlayer struct{}

func NewMusicAppPlayer() *MusicAppPlayer {
	return &MusicAppPlayer{}
}

func (m *MusicAppPlayer) Status(_ context.Context) (Status, error) {
	return Status{}, ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) Play(_ context.Context) error {
	return ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) Pause(_ context.Context) error {
	return ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) TogglePlayPause(_ context.Context) error {
	return ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) Next(_ context.Context) error {
	return ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) Previous(_ context.Context) error {
	return ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) Seek(_ context.Context, _ float64) error {
	return ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) SetVolume(_ context.Context, _ int) error {
	return ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) SetShuffle(_ context.Context, _ bool) error {
	return ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) SetRepeat(_ context.Context, _ RepeatMode) error {
	return ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) PlayTrack(_ context.Context, _, _ string) error {
	return ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) ListPlaylists(_ context.Context) ([]Playlist, error) {
	return nil, ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) ListTracks(_ context.Context, _ string, _, _ int) ([]Track, int, error) {
	return nil, 0, ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) Artwork(_ context.Context) ([]byte, string, error) {
	return nil, "", ErrMusicAppUnavailable
}

func (m *MusicAppPlayer) IsRunning(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "pgrep", "-x", "Music")
	return cmd.Run() == nil
}

func (m *MusicAppPlayer) Launch(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "open", "-ga", "Music")
	return cmd.Run()
}
