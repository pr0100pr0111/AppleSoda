package ui

import (
	"math/rand/v2"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/player"
)

func (m *Model) resolvePlayingContext() ([]player.Track, string) {
	if len(m.playingTracks) > 0 {
		return m.playingTracks, m.playingPLID
	}
	displayTracks := m.currentDisplayTracks()
	if len(displayTracks) > 0 {
		m.playingTracks = displayTracks
		m.playingPLID = m.currentPLID
		return m.playingTracks, m.playingPLID
	}
	if len(m.tracks) > 0 {
		m.playingTracks = m.tracks
		m.playingPLID = m.currentPLID
		return m.playingTracks, m.playingPLID
	}
	return nil, ""
}

func (m Model) findTrackIndex(tracks []player.Track, id string) int {
	if id == "" {
		return -1
	}
	for i, t := range tracks {
		if t.PersistentID == id {
			return i
		}
	}
	return -1
}

func (m Model) isTrackInPlayingList(id string) bool {
	if id == "" {
		return false
	}
	for _, t := range m.playingTracks {
		if t.PersistentID == id {
			return true
		}
	}
	return false
}

func (m *Model) initShufflePool(currentID string, tracks []player.Track) {
	m.shuffleHistory = nil
	if currentID != "" {
		m.shuffleHistory = []string{currentID}
	}
	m.shufflePool = make([]string, 0, len(tracks))
	for _, t := range tracks {
		if t.PersistentID != currentID {
			m.shufflePool = append(m.shufflePool, t.PersistentID)
		}
	}
	rand.Shuffle(len(m.shufflePool), func(i, j int) {
		m.shufflePool[i], m.shufflePool[j] = m.shufflePool[j], m.shufflePool[i]
	})
}

func (m *Model) syncSelectedTrackWithID(id string) {
	if m.playingPLID == m.currentPLID {
		displayTracks := m.currentDisplayTracks()
		for i, t := range displayTracks {
			if t.PersistentID == id {
				m.selectedTrack = i
				return
			}
		}
	}
}

func (m Model) shouldAutoAdvance(newStatus player.Status) bool {
	if m.playingTrackID == "" || m.playlistFinished {
		return false
	}
	if time.Since(m.lastPlayActionAt) < 1500*time.Millisecond {
		return false
	}
	if newStatus.State == player.StateStopped {
		return true
	}
	if newStatus.Track.PersistentID != "" && newStatus.Track.PersistentID != m.playingTrackID {
		if !m.isTrackInPlayingList(newStatus.Track.PersistentID) {
			return true
		}
	}
	return false
}

func (m Model) handleNextTrack() (tea.Model, tea.Cmd) {
	tracks, plID := m.resolvePlayingContext()
	if len(tracks) == 0 {
		return m, m.nextTrack()
	}

	m.lastPlayActionAt = time.Now()

	if m.status.Shuffle {
		if len(m.shufflePool) == 0 {
			if m.status.Repeat == player.RepeatAll {
				m.initShufflePool(m.playingTrackID, tracks)
			} else {
				m.playlistFinished = true
				m.playingTrackID = ""
				return m, m.nextTrack()
			}
		}
		if len(m.shufflePool) > 0 {
			nextID := m.shufflePool[0]
			m.shufflePool = m.shufflePool[1:]
			m.shuffleHistory = append(m.shuffleHistory, nextID)
			m.playingTrackID = nextID
			m.playlistFinished = false
			m.syncSelectedTrackWithID(nextID)
			return m, m.playTrack(plID, nextID)
		}
	}

	currIdx := m.findTrackIndex(tracks, m.playingTrackID)
	if currIdx == -1 {
		currIdx = m.findTrackIndex(tracks, m.status.Track.PersistentID)
	}
	if currIdx == -1 {
		currIdx = m.selectedTrack - 1
		if currIdx < -1 {
			currIdx = -1
		}
	}

	nextIdx := currIdx + 1
	if nextIdx < len(tracks) {
		t := tracks[nextIdx]
		m.playingTrackID = t.PersistentID
		m.playlistFinished = false
		m.syncSelectedTrackWithID(t.PersistentID)
		return m, m.playTrack(plID, t.PersistentID)
	}

	switch m.status.Repeat {
	case player.RepeatAll:
		t := tracks[0]
		m.playingTrackID = t.PersistentID
		m.playlistFinished = false
		m.syncSelectedTrackWithID(t.PersistentID)
		return m, m.playTrack(plID, t.PersistentID)
	case player.RepeatOne:
		if currIdx >= 0 && currIdx < len(tracks) {
			t := tracks[currIdx]
			return m, m.playTrack(plID, t.PersistentID)
		}
	default:
		m.playlistFinished = true
		m.playingTrackID = ""
		return m, m.nextTrack()
	}

	return m, nil
}

func (m Model) handlePrevTrack() (tea.Model, tea.Cmd) {
	if m.interpPos > 3.0 {
		return m, m.seek(0)
	}

	tracks, plID := m.resolvePlayingContext()
	if len(tracks) == 0 {
		return m, m.prevTrack()
	}

	m.lastPlayActionAt = time.Now()

	if m.status.Shuffle {
		if len(m.shuffleHistory) > 1 {
			m.shuffleHistory = m.shuffleHistory[:len(m.shuffleHistory)-1]
			prevID := m.shuffleHistory[len(m.shuffleHistory)-1]
			m.playingTrackID = prevID
			m.playlistFinished = false
			m.syncSelectedTrackWithID(prevID)
			return m, m.playTrack(plID, prevID)
		}
		return m, m.seek(0)
	}

	currIdx := m.findTrackIndex(tracks, m.playingTrackID)
	if currIdx == -1 {
		currIdx = m.findTrackIndex(tracks, m.status.Track.PersistentID)
	}
	if currIdx == -1 {
		currIdx = m.selectedTrack
	}

	prevIdx := currIdx - 1
	if prevIdx >= 0 && prevIdx < len(tracks) {
		t := tracks[prevIdx]
		m.playingTrackID = t.PersistentID
		m.playlistFinished = false
		m.syncSelectedTrackWithID(t.PersistentID)
		return m, m.playTrack(plID, t.PersistentID)
	}

	if m.status.Repeat == player.RepeatAll && len(tracks) > 0 {
		prevIdx = len(tracks) - 1
		t := tracks[prevIdx]
		m.playingTrackID = t.PersistentID
		m.playlistFinished = false
		m.syncSelectedTrackWithID(t.PersistentID)
		return m, m.playTrack(plID, t.PersistentID)
	}

	return m, m.seek(0)
}
