package ui

import (
	"fmt"
	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/theme"
)

func truncateStr(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxWidth {
		return s
	}
	if maxWidth <= 1 {
		return "…"
	}
	return string(runes[:maxWidth-1]) + "…"
}

func formatDuration(seconds float64) string {
	if seconds <= 0 {
		return "0:00"
	}
	m := int(seconds) / 60
	s := int(seconds) % 60
	return fmt.Sprintf("%d:%02d", m, s)
}

func renderProgressBar(width int, ratio float64, th theme.Theme) string {
	if width <= 0 {
		return ""
	}
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	filled := int(float64(width) * ratio)
	if filled > width {
		filled = width
	}
	empty := width - filled

	filledStyle := lipgloss.NewStyle().Foreground(th.ProgressFilled)
	emptyStyle := lipgloss.NewStyle().Foreground(th.ProgressEmpty)

	return filledStyle.Render(strings.Repeat("━", filled)) +
		emptyStyle.Render(strings.Repeat("━", empty))
}

func renderVolumeBar(width int, volume int, th theme.Theme) string {
	if width <= 0 {
		return ""
	}
	ratio := float64(volume) / 100.0
	filled := int(float64(width) * ratio)
	if filled > width {
		filled = width
	}
	empty := width - filled

	filledStyle := lipgloss.NewStyle().Foreground(th.Accent)
	emptyStyle := lipgloss.NewStyle().Foreground(th.ProgressEmpty)

	return filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))
}

func padRight(s string, width int) string {
	runes := []rune(s)
	if len(runes) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(runes))
}
