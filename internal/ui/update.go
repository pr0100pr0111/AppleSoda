package ui

import (
	"context"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/config"
	"github.com/pr0100pr0111/AppleSoda/internal/library"
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
			const settle = 1500 * time.Millisecond

			savedInterpPos := m.interpPos
			savedVolume := m.status.Volume
			savedInterpState := m.interpState

			m.statusErr = ""
			m.status = msg.status
			m.lastPollAt = time.Now()

			if time.Since(m.lastSeekCmd) <= settle {
				m.interpPos = savedInterpPos
				m.status.Position = savedInterpPos
			} else {
				m.interpPos = msg.status.Position
			}

			if time.Since(m.lastVolumeCmd) <= settle {
				m.status.Volume = savedVolume
			}

			if time.Since(m.lastPauseCmd) <= settle {
				m.interpState = savedInterpState
			} else {
				m.interpState = msg.status.State
			}

			if m.playingTrackID != "" && msg.status.Track.PersistentID != "" && msg.status.Track.PersistentID != m.playingTrackID {
				if m.isTrackInPlayingList(msg.status.Track.PersistentID) {
					m.playingTrackID = msg.status.Track.PersistentID
				}
			}

			if m.shouldAutoAdvance(msg.status) {
				if m.status.Repeat == player.RepeatOne {
					m.lastPlayActionAt = time.Now()
					return m, m.playTrack(m.playingPLID, m.playingTrackID)
				}
				return m.handleNextTrack()
			}
		}
		return m, nil
	case playlistsMsg:
		if msg.err != nil {
			m.statusErr = msg.err.Error()
		} else {
			m.playlists = msg.playlists
			m.playlistsLoaded = true
			if m.store != nil {
				m.store.SetPlaylists(msg.playlists)
			}
			if len(m.playlists) > 0 && m.currentPLID == "" {
				m.currentPLID = m.playlists[0].PersistentID
				if m.store != nil {
					if cachedTrs, ok := m.store.GetTracks(m.currentPLID); ok && len(cachedTrs) > 0 {
						m.tracks = cachedTrs
						m.tracksTotal = len(cachedTrs)
						m.applyFilter()
					}
				}
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
			if m.store != nil {
				m.store.SetTracks(msg.playlistID, msg.tracks)
			}
			m.applyFilter()
		}
		return m, nil
	case actionDoneMsg:
		if msg.err != nil {
			m.statusErr = msg.err.Error()
		}
		return m, nil
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
		case "esc":
			m.filterMode = false
			m.filterText = ""
			m.applyFilter()
			return m, nil
		case "enter":
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
	case "esc":
		if m.filterText != "" {
			m.filterText = ""
			m.applyFilter()
			return m, nil
		}
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
				if m.store != nil {
					if cachedTrs, ok := m.store.GetTracks(pl.PersistentID); ok && len(cachedTrs) > 0 {
						m.tracks = cachedTrs
						m.tracksTotal = len(cachedTrs)
						m.currentPLID = pl.PersistentID
						m.selectedTrack = 0
						m.trackScroll = 0
						m.applyFilter()
					}
				}
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
		if m.interpState == player.StatePlaying {
			m.interpState = player.StatePaused
		} else {
			m.interpState = player.StatePlaying
		}
		m.lastPauseCmd = time.Now()
		return m, m.togglePlayPause()
	case "n":
		return m.handleNextTrack()
	case "p":
		return m.handlePrevTrack()
	case "z":
		m.scrollToPlaying()
		return m, nil
	case "[", ",", "shift+left", "х", "б":
		newPos := m.interpPos - 5
		if newPos < 0 {
			newPos = 0
		}
		m.interpPos = newPos
		m.lastSeekCmd = time.Now()
		return m, m.seek(newPos)
	case "]", ".", "shift+right", "ъ", "ю":
		newPos := m.interpPos + 5
		if m.status.Track.Duration > 0 && newPos > m.status.Track.Duration {
			newPos = m.status.Track.Duration
		}
		m.interpPos = newPos
		m.lastSeekCmd = time.Now()
		return m, m.seek(newPos)
	case "{", "<", "Х", "Б":
		newPos := m.interpPos - 30
		if newPos < 0 {
			newPos = 0
		}
		m.interpPos = newPos
		m.lastSeekCmd = time.Now()
		return m, m.seek(newPos)
	case "}", ">", "Ъ", "Ю":
		newPos := m.interpPos + 30
		if m.status.Track.Duration > 0 && newPos > m.status.Track.Duration {
			newPos = m.status.Track.Duration
		}
		m.interpPos = newPos
		m.lastSeekCmd = time.Now()
		return m, m.seek(newPos)
	case "+", "=":
		newVol := m.status.Volume + 5
		if newVol > 100 {
			newVol = 100
		}
		m.status.Volume = newVol
		m.lastVolumeCmd = time.Now()
		return m, m.setVolume(newVol)
	case "-":
		newVol := m.status.Volume - 5
		if newVol < 0 {
			newVol = 0
		}
		m.status.Volume = newVol
		m.lastVolumeCmd = time.Now()
		return m, m.setVolume(newVol)
	case "s":
		newShuffle := !m.status.Shuffle
		m.status.Shuffle = newShuffle
		if newShuffle {
			tracks, _ := m.resolvePlayingContext()
			m.initShufflePool(m.playingTrackID, tracks)
		}
		return m, m.setShuffle(newShuffle)
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
				if m.store != nil {
					if cachedTrs, ok := m.store.GetTracks(pl.PersistentID); ok && len(cachedTrs) > 0 {
						m.tracks = cachedTrs
						m.tracksTotal = len(cachedTrs)
						m.currentPLID = pl.PersistentID
						m.selectedTrack = 0
						m.trackScroll = 0
						m.applyFilter()
					}
				}
				m.tracksLoading = true
				return m, m.loadTracks(pl.PersistentID)
			}
		}
	case panelTracks:
		curTracks := m.currentDisplayTracks()
		if len(curTracks) > 0 && m.selectedTrack < len(curTracks) {
			t := curTracks[m.selectedTrack]
			m.playingPLID = m.currentPLID
			m.playingTracks = curTracks
			m.playingTrackID = t.PersistentID
			m.playlistFinished = false
			m.lastPlayActionAt = time.Now()
			if m.status.Shuffle {
				m.initShufflePool(t.PersistentID, curTracks)
			}
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
	m.filteredIdx = library.FilterIndices(m.tracks, m.filterText)
	m.selectedTrack = 0
	m.trackScroll = 0
}

func (m *Model) scrollToPlaying() {
	if !m.status.HasTrack {
		return
	}
	displayTracks := m.currentDisplayTracks()
	for i, t := range displayTracks {
		if t.PersistentID == m.status.Track.PersistentID {
			m.selectedTrack = i
			m.active = panelTracks
			return
		}
	}
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
