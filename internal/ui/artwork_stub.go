package ui

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/theme"
)

func renderArtworkStub(width, height int, th theme.Theme) string {
	if width < 8 || height < 4 {
		return ""
	}

	size := height
	if size*2 > width {
		size = width / 2
	}
	if size < 4 {
		size = 4
	}

	cols := size * 2
	rows := size

	lines := make([]string, rows)
	bgStyle := lipgloss.NewStyle().Foreground(th.Border).Background(th.Background)
	noteStyle := lipgloss.NewStyle().Foreground(th.Primary).Background(th.Background)
	accentStyle := lipgloss.NewStyle().Foreground(th.Secondary).Background(th.Background)

	for r := 0; r < rows; r++ {
		var sb strings.Builder
		for c := 0; c < cols; c++ {
			switch {
			case (r == 1 || r == 2) && (c >= cols/2-1 && c <= cols/2+3):
				sb.WriteString(noteStyle.Render("▀"))
			case (r >= 2 && r <= rows-3) && (c == cols/2-1 || c == cols/2+3):
				sb.WriteString(accentStyle.Render("█"))
			case (r == rows-2 || r == rows-3) && ((c >= cols/2-3 && c <= cols/2-1) || (c >= cols/2+1 && c <= cols/2+3)):
				sb.WriteString(noteStyle.Render("▄"))
			case r == 0 || r == rows-1 || c == 0 || c == cols-1:
				sb.WriteString(bgStyle.Render("·"))
			default:
				sb.WriteString(" ")
			}
		}
		lines[r] = sb.String()
	}

	return strings.Join(lines, "\n")
}
