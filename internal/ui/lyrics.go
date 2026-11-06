package ui

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/lyrics"
)

func (m Model) renderLyrics(width, height int) string {
	if width < 10 || height < 1 {
		return ""
	}

	if !m.status.HasTrack {
		return ""
	}

	if m.lyricsLoading {
		mutedStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
		return mutedStyle.Render(truncateStr("Loading lyrics...", width))
	}

	if m.lyricsResult.Synced && len(m.lyricsResult.Lines) > 0 {
		return m.renderSyncedLyrics(width, height)
	}

	if m.lyricsResult.Plain != "" {
		return m.renderPlainLyrics(width, height)
	}

	mutedStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
	return mutedStyle.Render(truncateStr("Текст не найден", width))
}

func (m Model) renderSyncedLyrics(width, height int) string {
	if height <= 0 || width <= 0 {
		return ""
	}
	lines := m.lyricsResult.Lines
	activeIdx := lyrics.ActiveIndex(lines, m.interpPos)

	activeStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.NowPlaying)
	mutedStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)

	start := 0
	if len(lines) > height {
		half := height / 2
		start = activeIdx - half
		if start < 0 {
			start = 0
		}
		if start+height > len(lines) {
			start = len(lines) - height
		}
	}

	end := start + height
	if end > len(lines) {
		end = len(lines)
	}

	var rendered []string
	for i := start; i < end; i++ {
		txt := truncateStr(lines[i].Text, width-2)
		if i == activeIdx {
			rendered = append(rendered, activeStyle.Render("▶ "+txt))
		} else {
			rendered = append(rendered, mutedStyle.Render("  "+txt))
		}
	}
	for len(rendered) < height {
		rendered = append(rendered, "")
	}

	return strings.Join(rendered, "\n")
}

func (m Model) renderPlainLyrics(width, height int) string {
	rawLines := strings.Split(m.lyricsResult.Plain, "\n")
	var cleaned []string
	for _, l := range rawLines {
		cleaned = append(cleaned, strings.TrimRight(l, "\r"))
	}

	start := 0
	if len(cleaned) > height && m.status.Track.Duration > 0 {
		ratio := m.interpPos / m.status.Track.Duration
		if ratio > 1 {
			ratio = 1
		} else if ratio < 0 {
			ratio = 0
		}
		maxStart := len(cleaned) - height
		start = int(float64(maxStart) * ratio)
		if start < 0 {
			start = 0
		}
	}

	end := start + height
	if end > len(cleaned) {
		end = len(cleaned)
	}

	lineStyle := lipgloss.NewStyle().Foreground(m.theme.Foreground)
	var rendered []string
	for i := start; i < end; i++ {
		rendered = append(rendered, lineStyle.Render(truncateStr(cleaned[i], width)))
	}
	for len(rendered) < height {
		rendered = append(rendered, "")
	}

	return strings.Join(rendered, "\n")
}
