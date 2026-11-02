package ui

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/theme"
)

func renderHelpModal(width, height int, th theme.Theme) string {
	shortcuts := [][2]string{
		{"Tab / Shift+Tab", "Next / prev panel"},
		{"j, k / ↓, ↑", "Move cursor down / up"},
		{"h, l / ←, →", "Navigate panels (left / right)"},
		{"g / G", "First / last item"},
		{"Ctrl-d / Ctrl-u", "Page down / up (10 items)"},
		{"Enter", "Play selected item"},
		{"Space", "Play / pause"},
		{"n / p", "Next / previous track"},
		{"[ / ], Shift+←/→", "Seek -5s / +5s"},
		{", / .", "Seek -5s / +5s (alternative)"},
		{"{ / }, < / >", "Seek -30s / +30s"},
		{"+ / -", "Volume up / down"},
		{"s", "Toggle shuffle (random order)"},
		{"r", "Cycle repeat (off/all/one)"},
		{"/", "Filter tracklist (Esc to cancel)"},
		{"z", "Jump to currently playing track"},
		{"v", "Toggle visualizer mode (bars/mirror)"},
		{"t", "Cycle theme (apple/catppuccin/etc)"},
		{"c / L", "Toggle cover / toggle lyrics"},
		{"o", "Launch Music.app"},
		{"R", "Reload library"},
		{"?", "Close this help window"},
		{"q / Q", "Quit (leave playing) / Quit & pause"},
	}

	keyStyle := lipgloss.NewStyle().Bold(true).Foreground(th.Primary).Width(20)
	descStyle := lipgloss.NewStyle().Foreground(th.Foreground)

	var rows []string
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(th.Accent).MarginBottom(1)
	rows = append(rows, titleStyle.Render("AppleSoda Keyboard Shortcuts"))

	for _, s := range shortcuts {
		rows = append(rows, keyStyle.Render(s[0])+" "+descStyle.Render(s[1]))
	}

	content := strings.Join(rows, "\n")
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(th.Primary).
		Background(th.Background).
		Padding(1, 2)

	renderedBox := boxStyle.Render(content)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, renderedBox)
}
