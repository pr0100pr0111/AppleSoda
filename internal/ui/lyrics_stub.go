package ui

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/theme"
)

type demoLyricLine struct {
	timeSec float64
	text    string
}

var sampleLyrics = []demoLyricLine{
	{0, "♪ (Intro instrumental) ♪"},
	{8, "Driving down the neon highway"},
	{16, "City lights reflect the rain"},
	{24, "Shadows dance across the dashboard"},
	{32, "Echoes of a forgotten name"},
	{40, "Turn the radio up higher"},
	{48, "Feel the rhythm in the night"},
	{56, "We are chasing after midnight"},
	{64, "Lost inside the synthwave light"},
	{72, "♪ (Synth solo) ♪"},
	{88, "Endless summer never fades away"},
	{96, "Memories frozen in the sound"},
	{104, "Catch the frequency of yesterday"},
	{112, "Never looking back down"},
	{120, "♪ (Outro) ♪"},
}

func renderLyricsStub(width, height int, currentSec float64, th theme.Theme) string {
	if width < 10 || height < 3 {
		return ""
	}

	activeIdx := 0
	for i, l := range sampleLyrics {
		if currentSec >= l.timeSec {
			activeIdx = i
		} else {
			break
		}
	}

	activeStyle := lipgloss.NewStyle().Bold(true).Foreground(th.NowPlaying)
	mutedStyle := lipgloss.NewStyle().Foreground(th.Muted)

	half := height / 2
	start := activeIdx - half
	end := start + height

	var lines []string
	for i := start; i < end; i++ {
		if i < 0 || i >= len(sampleLyrics) {
			lines = append(lines, "")
			continue
		}
		txt := truncateStr(sampleLyrics[i].text, width-2)
		if i == activeIdx {
			lines = append(lines, activeStyle.Render("▶ "+txt))
		} else {
			lines = append(lines, mutedStyle.Render("  "+txt))
		}
	}

	return strings.Join(lines, "\n")
}
