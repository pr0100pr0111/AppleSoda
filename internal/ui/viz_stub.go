package ui

import (
	"math"
	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/pr0100pr0111/AppleSoda/internal/theme"
)

var subcellBlocks = []rune{' ', ' ', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

func renderSimulatedVisualizer(width, height int, mode string, phase float64, isPlaying bool, th theme.Theme) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	numBars := width / 2
	if numBars <= 0 {
		numBars = 1
	}

	levels := make([]float64, numBars)
	for i := 0; i < numBars; i++ {
		if !isPlaying {
			levels[i] = 0.05
			continue
		}
		fi := float64(i)
		w1 := math.Sin(fi*0.35 + phase*1.8)
		w2 := math.Cos(fi*0.2 - phase*1.2)
		w3 := math.Sin(fi*0.8+phase*3.0) * 0.5
		raw := (w1 + w2 + w3 + 2.5) / 5.0
		if raw < 0 {
			raw = 0
		}
		if raw > 1 {
			raw = 1
		}
		factor := 1.0 - math.Pow(float64(i)/float64(numBars), 1.8)*0.5
		levels[i] = raw * factor
	}

	lines := make([]string, height)
	gradLen := len(th.VizGradient)
	if gradLen == 0 {
		gradLen = 1
	}

	if mode == "mirror" {
		midRow := height / 2
		for row := 0; row < height; row++ {
			var b strings.Builder
			distFromMid := math.Abs(float64(row - midRow))
			maxDist := float64(midRow)
			if maxDist == 0 {
				maxDist = 1
			}
			for i := 0; i < numBars; i++ {
				colorIdx := int(float64(i) / float64(numBars) * float64(gradLen))
				if colorIdx >= gradLen {
					colorIdx = gradLen - 1
				}
				style := lipgloss.NewStyle().Foreground(th.VizGradient[colorIdx])

				barMaxDist := levels[i] * maxDist
				if distFromMid <= barMaxDist {
					b.WriteString(style.Render("█ "))
				} else if distFromMid-barMaxDist < 0.6 {
					b.WriteString(style.Render("▄ "))
				} else {
					b.WriteString("  ")
				}
			}
			lines[row] = b.String()
		}
	} else {
		for row := 0; row < height; row++ {
			var b strings.Builder
			rowFromBottom := height - 1 - row
			for i := 0; i < numBars; i++ {
				colorIdx := int(float64(i) / float64(numBars) * float64(gradLen))
				if colorIdx >= gradLen {
					colorIdx = gradLen - 1
				}
				style := lipgloss.NewStyle().Foreground(th.VizGradient[colorIdx])

				barUnits := levels[i] * float64(height)
				fullBlocks := int(barUnits)
				if rowFromBottom < fullBlocks {
					b.WriteString(style.Render("█ "))
				} else if rowFromBottom == fullBlocks {
					fraction := barUnits - float64(fullBlocks)
					subIdx := int(fraction * float64(len(subcellBlocks)-1))
					if subIdx >= len(subcellBlocks) {
						subIdx = len(subcellBlocks) - 1
					}
					b.WriteString(style.Render(string(subcellBlocks[subIdx]) + " "))
				} else {
					b.WriteString("  ")
				}
			}
			lines[row] = b.String()
		}
	}

	return strings.Join(lines, "\n")
}
