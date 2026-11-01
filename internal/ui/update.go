package ui

import (
	"context"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/config"
	"github.com/pr0100pr0111/AppleSoda/internal/player"
	"github.com/pr0100pr0111/AppleSoda/internal/theme"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	case tea.MouseClickMsg:
		return m.handleMouseClick(msg)
	case tea.MouseWheelMsg:
		return m.handleMouseWheel(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.tooSmall = m.width < minWidth || m.height < minHeight
		return m, nil
	case statusMsg:
		if msg.err != nil {
			m.statusErr = msg.err.Error()
		} else {
			m.statusErr = ""
			m.status = msg.status
			m.interpPos = msg.status.Position
			m.interpState = msg.status.State
			m.lastPollAt = time.Now()
		}
		return m, nil
	case playlistsMsg:
		if msg.err != nil {
			m.statusErr = msg.err.Error()
		} else {
			m.playlists = msg.playlists
			m.playlistsLoaded = true
			if len(m.playlists) > 0 && m.currentPLID == "" {
				m.currentPLID = m.playlists[0].PersistentID
				m.tracksLoading = true
				return m, m.loadTracks(m.currentPLID)
			}
		}
		return m, nil
	case tracksMsg:
		m.tracksLoading = false
		if msg.err != nil {
			m.statusErr = msg.err.Error()
		} else {
			m.tracks = msg.tracks
			m.tracksTotal = msg.total
			m.currentPLID = msg.playlistID
			m.selectedTrack = 0
			m.trackScroll = 0
			m.applyFilter()
		}
		return m, nil
	case actionDoneMsg:
		if msg.err != nil {
			m.statusErr = msg.err.Error()
		}
		return m, m.pollStatus()
	case tickMsg:
		m.vizPhase += 0.1
		if m.interpState == player.StatePlaying {
			m.interpPos += 0.033
			if m.status.Track.Duration > 0 && m.interpPos > m.status.Track.Duration {
				m.interpPos = m.status.Track.Duration
			}
		}
		return m, m.tickCmd()
	case pollMsg:
		return m, tea.Batch(m.pollStatus(), m.pollCmd())
	}
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if m.showHelp {
		if key == "?" || key == "esc" || key == "q" || key == "enter" {
			m.showHelp = false
			return m, nil
		}
		return m, nil
	}

	if m.filterMode {
		switch key {
		case "esc", "enter":
			m.filterMode = false
			return m, nil
		case "backspace":
			runes := []rune(m.filterText)
			if len(runes) > 0 {
				m.filterText = string(runes[:len(runes)-1])
				m.applyFilter()
			}
			return m, nil
		default:
			if len(msg.Text) > 0 {
				m.filterText += msg.Text
				m.applyFilter()
			}
			return m, nil
		}
	}

	switch key {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit
	case "Q":
		m.quitting = true
		return m, tea.Sequence(m.playerPauseCmd(), tea.Quit)
	case "?":
		m.showHelp = true
		return m, nil
	case "tab":
		m.active = (m.active + 1) % 3
		return m, nil
	case "shift+tab":
		m.active = (m.active + 2) % 3
		return m, nil
	case "j", "down":
		m.moveCursor(1)
		return m, nil
	case "k", "up":
		m.moveCursor(-1)
		return m, nil
	case "g", "home":
		m.moveCursorStart()
		return m, nil
	case "G", "end":
		m.moveCursorEnd()
		return m, nil
	case "ctrl+d":
		m.moveCursor(10)
		return m, nil
	case "ctrl+u":
		m.moveCursor(-10)
		return m, nil
	case "h", "left":
		if m.active == panelNowPlaying {
			m.active = panelTracks
		} else if m.active == panelTracks {
			m.active = panelLibrary
		}
		return m, nil
	case "l", "right":
		if m.active == panelLibrary && len(m.playlists) > 0 {
			m.active = panelTracks
			pl := m.playlists[m.selectedPL]
			if pl.PersistentID != m.currentPLID {
				m.tracksLoading = true
				return m, m.loadTracks(pl.PersistentID)
			}
		} else if m.active == panelTracks {
			m.active = panelNowPlaying
		}
		return m, nil
	case "enter":
		return m.handleEnter()
	case "space":
		return m, m.togglePlayPause()
	case "n":
		return m, m.nextTrack()
	case "p":
		return m, m.prevTrack()
	case "[", ",", "shift+left", "х", "б":
		return m, m.seek(m.interpPos - 5)
	case "]", ".", "shift+right", "ъ", "ю":
		return m, m.seek(m.interpPos + 5)
	case "{", "<", "Х", "Б":
		return m, m.seek(m.interpPos - 30)
	case "}", ">", "Ъ", "Ю":
		return m, m.seek(m.interpPos + 30)
	case "+", "=":
		return m, m.setVolume(m.status.Volume + 5)
	case "-":
		return m, m.setVolume(m.status.Volume - 5)
	case "s":
		return m, m.setShuffle(!m.status.Shuffle)
	case "r":
		return m, m.setRepeat(m.status.Repeat.Next())
	case "/":
		m.filterMode = true
		m.filterText = ""
		m.applyFilter()
		return m, nil
	case "v":
		if m.vizMode == "bars" {
			m.vizMode = "mirror"
		} else {
			m.vizMode = "bars"
		}
		return m, nil
	case "t":
		nextName := theme.NextTheme(m.theme.Name)
		m.theme = theme.Get(nextName)
		_ = config.SaveState(config.State{Theme: nextName})
		return m, nil
	case "c":
		m.showCover = !m.showCover
		return m, nil
	case "L":
		m.showLyrics = !m.showLyrics
		return m, nil
	case "o":
		return m, m.launchMusic()
	case "R":
		return m, tea.Batch(m.loadPlaylists(), m.pollStatus())
	}

	return m, nil
}

func (m Model) playerPauseCmd() tea.Cmd {
	return func() tea.Msg {
		_ = m.player.Pause(context.Background())
		return nil
	}
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.active {
	case panelLibrary:
		if len(m.playlists) > 0 {
			pl := m.playlists[m.selectedPL]
			m.active = panelTracks
			if pl.PersistentID != m.currentPLID {
				m.tracksLoading = true
				return m, m.loadTracks(pl.PersistentID)
			}
		}
	case panelTracks:
		curTracks := m.currentDisplayTracks()
		if len(curTracks) > 0 && m.selectedTrack < len(curTracks) {
			t := curTracks[m.selectedTrack]
			return m, m.playTrack(m.currentPLID, t.PersistentID)
		}
	case panelNowPlaying:
		return m, m.togglePlayPause()
	}
	return m, nil
}

func (m *Model) moveCursor(delta int) {
	switch m.active {
	case panelLibrary:
		total := len(m.playlists)
		if total == 0 {
			return
		}
		m.selectedPL += delta
		if m.selectedPL < 0 {
			m.selectedPL = 0
		}
		if m.selectedPL >= total {
			m.selectedPL = total - 1
		}
	case panelTracks:
		total := len(m.currentDisplayTracks())
		if total == 0 {
			return
		}
		m.selectedTrack += delta
		if m.selectedTrack < 0 {
			m.selectedTrack = 0
		}
		if m.selectedTrack >= total {
			m.selectedTrack = total - 1
		}
	}
}

func (m *Model) moveCursorStart() {
	switch m.active {
	case panelLibrary:
		m.selectedPL = 0
		m.plScroll = 0
	case panelTracks:
		m.selectedTrack = 0
		m.trackScroll = 0
	}
}

func (m *Model) moveCursorEnd() {
	switch m.active {
	case panelLibrary:
		if len(m.playlists) > 0 {
			m.selectedPL = len(m.playlists) - 1
		}
	case panelTracks:
		total := len(m.currentDisplayTracks())
		if total > 0 {
			m.selectedTrack = total - 1
		}
	}
}

func (m *Model) applyFilter() {
	if m.filterText == "" {
		m.filteredIdx = nil
		return
	}
	query := strings.ToLower(m.filterText)
	var matches []int
	for i, t := range m.tracks {
		name := strings.ToLower(t.Name)
		artist := strings.ToLower(t.Artist)
		album := strings.ToLower(t.Album)
		if strings.Contains(name, query) || strings.Contains(artist, query) || strings.Contains(album, query) {
			matches = append(matches, i)
		}
	}
	m.filteredIdx = matches
	m.selectedTrack = 0
	m.trackScroll = 0
}

func (m Model) currentDisplayTracks() []player.Track {
	if m.filterText != "" && m.filteredIdx != nil {
		res := make([]player.Track, 0, len(m.filteredIdx))
		for _, idx := range m.filteredIdx {
			if idx < len(m.tracks) {
				res = append(res, m.tracks[idx])
			}
		}
		return res
	}
	return m.tracks
}

func (m Model) handleMouseClick(msg tea.MouseClickMsg) (tea.Model, tea.Cmd) {
	if !m.cfg.UI.Mouse {
		return m, nil
	}
	mTarget := msg.Mouse()
	if mTarget.Button != tea.MouseLeft {
		return m, nil
	}

	isDouble := time.Since(m.lastClickTime) < 400*time.Millisecond && m.lastClickRow == mTarget.Y && m.lastClickPan == m.active
	m.lastClickTime = time.Now()
	m.lastClickRow = mTarget.Y
	m.lastClickPan = m.active

	leftColWidth := 30
	if m.width >= 120 {
		leftColWidth = m.width / 4
		if leftColWidth < 25 {
			leftColWidth = 25
		}
	}

	if mTarget.X < leftColWidth && mTarget.Y > 1 && mTarget.Y < m.height-8 {
		m.active = panelLibrary
		clickedRow := mTarget.Y - 2 + m.plScroll
		if clickedRow >= 0 && clickedRow < len(m.playlists) {
			m.selectedPL = clickedRow
			if isDouble {
				return m.handleEnter()
			}
		}
		return m, nil
	}

	tracksWidth := m.width / 2
	if mTarget.X >= leftColWidth && mTarget.X < leftColWidth+tracksWidth && mTarget.Y > 1 && mTarget.Y < m.height-8 {
		m.active = panelTracks
		clickedRow := mTarget.Y - 2 + m.trackScroll
		displayTracks := m.currentDisplayTracks()
		if clickedRow >= 0 && clickedRow < len(displayTracks) {
			m.selectedTrack = clickedRow
			if isDouble {
				return m.handleEnter()
			}
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleMouseWheel(msg tea.MouseWheelMsg) (tea.Model, tea.Cmd) {
	if !m.cfg.UI.Mouse {
		return m, nil
	}
	mTarget := msg.Mouse()
	if mTarget.Button == tea.MouseWheelUp {
		m.moveCursor(-2)
	} else if mTarget.Button == tea.MouseWheelDown {
		m.moveCursor(2)
	}
	return m, nil
}
