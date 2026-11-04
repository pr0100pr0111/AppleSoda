package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/pr0100pr0111/AppleSoda/internal/player"
)

func (m Model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("")
	}

	if m.tooSmall {
		content := fmt.Sprintf("Terminal window is too small (%dx%d).\nPlease resize to at least %dx%d.", m.width, m.height, minWidth, minHeight)
		styled := lipgloss.NewStyle().
			Foreground(m.theme.Primary).
			Bold(true).
			Render(content)
		placed := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, styled)
		v := tea.NewView(placed)
		v.AltScreen = true
		return v
	}

	if m.showHelp {
		placed := renderHelpModal(m.width, m.height, m.theme)
		v := tea.NewView(placed)
		v.AltScreen = true
		return v
	}

	content := m.renderFullLayout()
	v := tea.NewView(content)
	v.AltScreen = true
	if m.cfg.UI.Mouse {
		v.MouseMode = tea.MouseModeCellMotion
	}
	return v
}

func (m Model) renderFullLayout() string {
	headerHeight := 1
	transportHeight := 3
	vizHeight := 7
	if m.height < 25 {
		vizHeight = 4
	}

	mainHeight := m.height - headerHeight - transportHeight - vizHeight
	if mainHeight < 6 {
		mainHeight = 6
	}

	header := m.renderHeader()
	mainArea := m.renderMainArea(mainHeight)
	vizArea := m.renderVisualizerArea(vizHeight)
	transport := m.renderTransportArea()

	return lipgloss.JoinVertical(lipgloss.Left, header, mainArea, vizArea, transport)
}

func (m Model) renderHeader() string {
	logoStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary)
	logo := logoStyle.Render("AppleSoda")

	themeStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
	thInfo := themeStyle.Render(fmt.Sprintf("[%s]", m.theme.Name))

	var statusTxt string
	if m.filterMode {
		filterStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Accent)
		statusTxt = filterStyle.Render(fmt.Sprintf(" / %s█", m.filterText))
	} else if m.filterText != "" {
		filterStyle := lipgloss.NewStyle().Foreground(m.theme.Accent)
		statusTxt = filterStyle.Render(fmt.Sprintf(" [filter: %s]", m.filterText))
	} else if m.statusErr != "" {
		errStyle := lipgloss.NewStyle().Foreground(m.theme.Primary)
		statusTxt = errStyle.Render(" " + truncateStr(m.statusErr, m.width-30))
	} else {
		hintStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
		statusTxt = hintStyle.Render(" Tab:switch  ?:help  q:quit")
	}

	left := logo + " " + thInfo + " " + statusTxt
	return lipgloss.NewStyle().Width(m.width).Render(left)
}

func (m Model) renderMainArea(height int) string {
	if m.width >= 120 {
		return m.renderWideMain(height)
	} else if m.width >= 80 {
		return m.renderMediumMain(height)
	}
	return m.renderNarrowMain(height)
}

func (m Model) renderWideMain(height int) string {
	leftWidth := m.width / 4
	if leftWidth < 22 {
		leftWidth = 22
	}
	rightWidth := m.width / 3
	if rightWidth < 28 {
		rightWidth = 28
	}
	centerWidth := m.width - leftWidth - rightWidth

	left := m.renderLibraryPanel(leftWidth, height)
	center := m.renderTracksPanel(centerWidth, height)
	right := m.renderNowPlayingPanel(rightWidth, height)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, center, right)
}

func (m Model) renderMediumMain(height int) string {
	rightWidth := m.width * 2 / 5
	leftWidth := m.width - rightWidth

	var left string
	if m.active == panelLibrary {
		left = m.renderLibraryPanel(leftWidth, height)
	} else {
		left = m.renderTracksPanel(leftWidth, height)
	}
	right := m.renderNowPlayingPanel(rightWidth, height)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m Model) renderNarrowMain(height int) string {
	switch m.active {
	case panelLibrary:
		return m.renderLibraryPanel(m.width, height)
	case panelTracks:
		return m.renderTracksPanel(m.width, height)
	default:
		return m.renderNowPlayingPanel(m.width, height)
	}
}

func (m Model) panelBorderStyle(active bool) lipgloss.Style {
	st := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder())
	if active {
		st = st.BorderForeground(m.theme.Primary)
	} else {
		st = st.BorderForeground(m.theme.Border)
	}
	return st
}

func (m Model) renderLibraryPanel(width, height int) string {
	style := m.panelBorderStyle(m.active == panelLibrary).
		Width(width).
		Height(height)

	innerW := width - 2
	innerH := height - 2
	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Foreground)
	lines := []string{headerStyle.Render("Playlists")}

	if len(m.playlists) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
		lines = append(lines, emptyStyle.Render("  No playlists"))
	} else {
		maxItems := innerH - 1
		if maxItems < 1 {
			maxItems = 1
		}

		start := m.selectedPL - maxItems/2
		if start < 0 {
			start = 0
		}
		end := start + maxItems
		if end > len(m.playlists) {
			end = len(m.playlists)
			start = end - maxItems
			if start < 0 {
				start = 0
			}
		}

		for i := start; i < end; i++ {
			pl := m.playlists[i]
			name := truncateStr(pl.Name, innerW-4)
			var item string
			if i == m.selectedPL && m.active == panelLibrary {
				itemStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary)
				item = itemStyle.Render("▶ " + name)
			} else if pl.PersistentID == m.currentPLID {
				itemStyle := lipgloss.NewStyle().Foreground(m.theme.Accent)
				item = itemStyle.Render("• " + name)
			} else {
				itemStyle := lipgloss.NewStyle().Foreground(m.theme.Foreground)
				item = itemStyle.Render("  " + name)
			}
			lines = append(lines, item)
		}
	}

	for len(lines) < innerH {
		lines = append(lines, "")
	}

	return style.Render(strings.Join(lines[:innerH], "\n"))
}

func (m Model) renderTracksPanel(width, height int) string {
	style := m.panelBorderStyle(m.active == panelTracks).
		Width(width).
		Height(height)

	innerW := width - 2
	innerH := height - 2
	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}

	curPlaylistName := "Tracks"
	for _, pl := range m.playlists {
		if pl.PersistentID == m.currentPLID {
			curPlaylistName = pl.Name
			break
		}
	}

	displayTracks := m.currentDisplayTracks()
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Foreground)
	headerText := fmt.Sprintf("%s (%d)", curPlaylistName, len(displayTracks))
	lines := []string{headerStyle.Render(truncateStr(headerText, innerW))}

	if m.tracksLoading {
		lines = append(lines, lipgloss.NewStyle().Foreground(m.theme.Muted).Render("  Loading tracks..."))
	} else if len(displayTracks) == 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(m.theme.Muted).Render("  No tracks found"))
	} else {
		maxItems := innerH - 1
		if maxItems < 1 {
			maxItems = 1
		}

		start := m.selectedTrack - maxItems/2
		if start < 0 {
			start = 0
		}
		end := start + maxItems
		if end > len(displayTracks) {
			end = len(displayTracks)
			start = end - maxItems
			if start < 0 {
				start = 0
			}
		}

		for i := start; i < end; i++ {
			t := displayTracks[i]
			durStr := formatDuration(t.Duration)
			availableNameW := innerW - len(durStr) - 6
			if availableNameW < 10 {
				availableNameW = 10
			}

			trackTitle := truncateStr(t.Name, availableNameW*3/5)
			trackArtist := truncateStr(t.Artist, availableNameW*2/5)
			combined := fmt.Sprintf("%s - %s", trackTitle, trackArtist)
			combined = padRight(combined, availableNameW)

			isPlayingTrack := (m.status.HasTrack && m.status.Track.PersistentID == t.PersistentID)
			isSelected := (i == m.selectedTrack && m.active == panelTracks)

			var line string
			if isPlayingTrack && isSelected {
				st := lipgloss.NewStyle().Bold(true).Foreground(m.theme.NowPlaying)
				line = st.Render("♫ " + combined + " " + durStr)
			} else if isPlayingTrack {
				st := lipgloss.NewStyle().Foreground(m.theme.NowPlaying)
				line = st.Render("♫ " + combined + " " + durStr)
			} else if isSelected {
				st := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary)
				line = st.Render("▶ " + combined + " " + durStr)
			} else {
				st := lipgloss.NewStyle().Foreground(m.theme.Foreground)
				durSt := lipgloss.NewStyle().Foreground(m.theme.Muted)
				line = "  " + st.Render(combined) + " " + durSt.Render(durStr)
			}
			lines = append(lines, line)
		}
	}

	for len(lines) < innerH {
		lines = append(lines, "")
	}

	return style.Render(strings.Join(lines[:innerH], "\n"))
}

func (m Model) renderNowPlayingPanel(width, height int) string {
	style := m.panelBorderStyle(m.active == panelNowPlaying).
		Width(width).
		Height(height)

	innerW := width - 2
	innerH := height - 2
	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}

	var sections []string

	if m.showCover && innerH >= 8 {
		artH := innerH / 3
		if artH > 8 {
			artH = 8
		}
		pid := m.status.Track.PersistentID
		art, ok := m.artworkCache[pid]
		if !ok {
			art = renderArtworkStub(innerW, artH, m.theme)
		}
		if art != "" {
			sections = append(sections, art)
		}
	}

	trackName := "Nothing playing"
	artistName := "Select a track and press Enter"
	albumName := ""
	if m.status.HasTrack {
		trackName = m.status.Track.Name
		artistName = m.status.Track.Artist
		albumName = m.status.Track.Album
	}

	titleSt := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Foreground)
	artistSt := lipgloss.NewStyle().Foreground(m.theme.Accent)
	albumSt := lipgloss.NewStyle().Foreground(m.theme.Muted)

	info := titleSt.Render(truncateStr(trackName, innerW)) + "\n" +
		artistSt.Render(truncateStr(artistName, innerW))
	if albumName != "" {
		info += "\n" + albumSt.Render(truncateStr(albumName, innerW))
	}
	sections = append(sections, info)

	if m.showLyrics {
		usedLines := 0
		for _, s := range sections {
			usedLines += strings.Count(s, "\n") + 1
		}
		remainingH := innerH - usedLines - 1
		if remainingH >= 3 {
			lyricView := renderLyricsStub(innerW, remainingH, m.interpPos, m.theme)
			sections = append(sections, lyricView)
		}
	}

	full := strings.Join(sections, "\n\n")
	allLines := strings.Split(full, "\n")
	if len(allLines) > innerH {
		allLines = allLines[:innerH]
	}
	for len(allLines) < innerH {
		allLines = append(allLines, "")
	}

	return style.Render(strings.Join(allLines, "\n"))
}

func (m Model) renderVisualizerArea(height int) string {
	borderColor := m.theme.Border
	if m.cfg.Visualizer.Pulse && m.interpState == player.StatePlaying {
		borderColor = m.theme.Accent
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(m.width).
		Height(height)

	innerW := m.width - 2
	innerH := height - 2
	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}

	tagStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
	modeStyle := lipgloss.NewStyle().Foreground(m.theme.Accent)
	title := fmt.Sprintf("Visualizer %s %s",
		tagStyle.Render("[simulated]"),
		modeStyle.Render(fmt.Sprintf("[%s]", m.vizMode)))

	bars := renderSimulatedVisualizer(innerW, innerH-1, m.vizMode, m.vizPhase, m.interpState == player.StatePlaying, m.theme)

	content := title + "\n" + bars
	return style.Render(content)
}

func (m Model) renderTransportArea() string {
	dur := m.status.Track.Duration
	cur := m.interpPos
	ratio := 0.0
	if dur > 0 {
		ratio = cur / dur
	}

	timeStr := fmt.Sprintf("%s / %s", formatDuration(cur), formatDuration(dur))

	var stateIcon string
	switch m.interpState {
	case player.StatePlaying:
		stateIcon = lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary).Render("▶ Play")
	case player.StatePaused:
		stateIcon = lipgloss.NewStyle().Bold(true).Foreground(m.theme.Accent).Render("⏸ Pause")
	default:
		stateIcon = lipgloss.NewStyle().Foreground(m.theme.Muted).Render("⏹ Stop")
	}

	volBarSize := 6
	if m.width < 80 {
		volBarSize = 4
	}
	volBar := renderVolumeBar(volBarSize, m.status.Volume, m.theme)
	volStr := fmt.Sprintf("Vol:%d%% [%s]", m.status.Volume, volBar)

	var shufStr string
	if m.status.Shuffle {
		shufStr = lipgloss.NewStyle().Bold(true).Foreground(m.theme.Accent).Render("[shuf:on]")
	} else {
		shufStr = lipgloss.NewStyle().Foreground(m.theme.Muted).Render("[shuf:off]")
	}

	var repStr string
	if m.status.Repeat != player.RepeatOff {
		repStr = lipgloss.NewStyle().Bold(true).Foreground(m.theme.Accent).Render(fmt.Sprintf("[rep:%s]", m.status.Repeat.String()))
	} else {
		repStr = lipgloss.NewStyle().Foreground(m.theme.Muted).Render("[rep:off]")
	}

	leftPart := fmt.Sprintf(" %s  %s ", stateIcon, timeStr)
	rightPart := fmt.Sprintf("  %s  %s %s", volStr, shufStr, repStr)

	leftW := ansi.StringWidth(leftPart)
	rightW := ansi.StringWidth(rightPart)

	barWidth := m.width - leftW - rightW
	if barWidth < 5 {
		barWidth = 5
	}

	progBar := renderProgressBar(barWidth, ratio, m.theme)
	line1 := leftPart + progBar + rightPart
	line1 = ansi.Truncate(line1, m.width, "")

	keysText := " Space:pause  n/p:skip  [/]:seek  +/-:vol  s:shuf  r:rep  v:viz  t:theme  c:art  L:lyr  ?:help"
	if m.width < 90 {
		keysText = " Space:pause  n/p:skip  [/]:seek  +/-:vol  s:shuf  ?:help"
	}
	keysLine := lipgloss.NewStyle().Foreground(m.theme.Muted).Render(truncateStr(keysText, m.width))

	return lipgloss.JoinVertical(lipgloss.Left, line1, keysLine)
}
