package ui

import (
	"context"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/artwork"
	"github.com/pr0100pr0111/AppleSoda/internal/config"
	"github.com/pr0100pr0111/AppleSoda/internal/library"
	"github.com/pr0100pr0111/AppleSoda/internal/lyrics"
	"github.com/pr0100pr0111/AppleSoda/internal/player"
	"github.com/pr0100pr0111/AppleSoda/internal/theme"
)

type panel int

const (
	panelLibrary panel = iota
	panelTracks
	panelNowPlaying
)

const (
	minWidth  = 60
	minHeight = 16
)

type statusMsg struct {
	status player.Status
	err    error
}

type playlistsMsg struct {
	playlists []player.Playlist
	err       error
}

type tracksMsg struct {
	tracks     []player.Track
	total      int
	playlistID string
	err        error
}

type actionDoneMsg struct {
	err error
}

type artworkMsg struct {
	persistentID string
	rendered     string
	err          error
}

type lyricsMsg struct {
	persistentID string
	result       lyrics.Result
	err          error
}

type tickMsg time.Time
type pollMsg time.Time

type Model struct {
	player    player.Player
	store     *library.Store
	cfg       config.Config
	theme     theme.Theme
	width     int
	height    int
	active    panel
	status    player.Status
	statusErr string

	playlists       []player.Playlist
	selectedPL      int
	plScroll        int
	playlistsLoaded bool

	tracks        []player.Track
	tracksTotal   int
	selectedTrack int
	trackScroll   int
	currentPLID   string
	tracksLoading bool

	filterMode  bool
	filterText  string
	filteredIdx []int

	showHelp   bool
	showCover  bool
	showLyrics bool
	vizMode    string
	vizPhase   float64
	tooSmall   bool
	quitting   bool

	lastPollAt    time.Time
	interpPos     float64
	interpState   player.PlayerState
	lastClickTime time.Time
	lastClickRow  int
	lastClickPan  panel

	lastVolumeCmd time.Time
	lastSeekCmd   time.Time
	lastPauseCmd  time.Time

	artworkCache   map[string]string
	artworkLoading bool
	lastArtworkID  string

	lyricsResult  lyrics.Result
	lyricsTrackID string
	lyricsLoading bool

	playingPLID      string
	playingTracks    []player.Track
	playingTrackID   string
	playlistFinished bool
	shufflePool      []string
	shuffleHistory   []string
	lastPlayActionAt time.Time
}

func NewModel(p player.Player, cfg config.Config, store *library.Store) Model {
	th := theme.Get(cfg.Theme)
	st := config.LoadState()
	if st.Theme != "" {
		th = theme.Get(st.Theme)
	}

	m := Model{
		player:       p,
		store:        store,
		cfg:          cfg,
		theme:        th,
		showCover:    cfg.UI.Cover,
		showLyrics:   cfg.UI.Lyrics,
		vizMode:      cfg.Visualizer.Mode,
		active:       panelLibrary,
		artworkCache: make(map[string]string),
	}

	if store != nil {
		cachedPls := store.GetPlaylists()
		if len(cachedPls) > 0 {
			m.playlists = cachedPls
			m.playlistsLoaded = true
			m.currentPLID = cachedPls[0].PersistentID
			if cachedTrs, ok := store.GetTracks(m.currentPLID); ok && len(cachedTrs) > 0 {
				m.tracks = cachedTrs
				m.tracksTotal = len(cachedTrs)
			}
		}
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.pollStatus(),
		m.loadPlaylists(),
		m.tickCmd(),
		m.pollCmd(),
	)
}

func (m Model) pollStatus() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		s, err := m.player.Status(ctx)
		return statusMsg{status: s, err: err}
	}
}

func (m Model) loadPlaylists() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		pl, err := m.player.ListPlaylists(ctx)
		return playlistsMsg{playlists: pl, err: err}
	}
}

func (m Model) loadTracks(playlistID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		tracks, total, err := m.player.ListTracks(ctx, playlistID, 0, 500)
		return tracksMsg{tracks: tracks, total: total, playlistID: playlistID, err: err}
	}
}

func (m Model) playTrack(playlistID string, trackID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := m.player.PlayTrack(ctx, playlistID, trackID)
		return actionDoneMsg{err: err}
	}
}

func (m Model) togglePlayPause() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := m.player.TogglePlayPause(ctx)
		return actionDoneMsg{err: err}
	}
}

func (m Model) nextTrack() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := m.player.Next(ctx)
		return actionDoneMsg{err: err}
	}
}

func (m Model) prevTrack() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := m.player.Previous(ctx)
		return actionDoneMsg{err: err}
	}
}

func (m Model) seek(pos float64) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := m.player.Seek(ctx, pos)
		return actionDoneMsg{err: err}
	}
}

func (m Model) setVolume(vol int) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := m.player.SetVolume(ctx, vol)
		return actionDoneMsg{err: err}
	}
}

func (m Model) setShuffle(on bool) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := m.player.SetShuffle(ctx, on)
		return actionDoneMsg{err: err}
	}
}

func (m Model) setRepeat(mode player.RepeatMode) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := m.player.SetRepeat(ctx, mode)
		return actionDoneMsg{err: err}
	}
}

func (m Model) launchMusic() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := m.player.Launch(ctx)
		return actionDoneMsg{err: err}
	}
}

func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(33*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) pollCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return pollMsg(t)
	})
}

func (m *Model) maybeLoadArtwork() tea.Cmd {
	if !m.showCover || !m.status.HasTrack || m.artworkLoading {
		return nil
	}
	if m.status.Track.PersistentID == m.lastArtworkID {
		return nil
	}
	artH := 4
	if m.height >= 30 {
		artH = 6
	} else if m.height < 20 {
		artH = 3
	}
	cols := artH * 2
	rows := artH
	m.artworkLoading = true
	return m.loadArtwork(cols, rows)
}

func (m Model) loadArtwork(cols, rows int) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		data, _, err := m.player.Artwork(ctx)
		if err != nil {
			return artworkMsg{err: err}
		}
		pid := m.status.Track.PersistentID
		rendered, err := artwork.Render(data, cols, rows, pid)
		return artworkMsg{persistentID: pid, rendered: rendered, err: err}
	}
}

func (m *Model) maybeLoadLyrics() tea.Cmd {
	if !m.showLyrics || !m.status.HasTrack || m.lyricsLoading {
		return nil
	}
	if m.status.Track.PersistentID == m.lyricsTrackID {
		return nil
	}
	m.lyricsResult = lyrics.Result{}
	m.lyricsLoading = true
	return m.loadLyrics()
}

func (m Model) loadLyrics() tea.Cmd {
	return func() tea.Msg {
		pid := m.status.Track.PersistentID
		artist := m.status.Track.Artist
		track := m.status.Track.Name
		album := m.status.Track.Album
		dur := m.status.Track.Duration

		if cached, ok := lyrics.Load(pid); ok {
			return lyricsMsg{persistentID: pid, result: cached}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		r, err := lyrics.Fetch(ctx, artist, track, album, dur)
		if (!r.Synced && r.Plain == "") || err != nil {
			ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel2()
			plain, _ := m.player.Lyrics(ctx2)
			if plain != "" {
				parsed := lyrics.ParseLRC(plain)
				if len(parsed) > 0 {
					r = lyrics.Result{Lines: parsed, Synced: true}
				} else {
					r = lyrics.Result{Plain: plain}
				}
			}
		}

		lyrics.Save(pid, r)
		return lyricsMsg{persistentID: pid, result: r, err: err}
	}
}
