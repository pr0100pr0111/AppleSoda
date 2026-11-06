package lyrics

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Line struct {
	TimeSec float64
	Text    string
}

func ParseLRC(src string) []Line {
	var lines []Line
	var offsetMs float64

	rawLines := strings.Split(src, "\n")

	for _, raw := range rawLines {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		if !strings.HasPrefix(raw, "[") {
			continue
		}

		var timestamps []float64
		rest := raw

		for strings.HasPrefix(rest, "[") {
			end := strings.Index(rest, "]")
			if end == -1 {
				break
			}
			tag := rest[1:end]
			rest = rest[end+1:]

			if strings.HasPrefix(tag, "offset:") {
				val := strings.TrimPrefix(tag, "offset:")
				ms, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
				if err == nil {
					offsetMs = ms
				}
				break
			}

			ts, err := parseTimestamp(tag)
			if err != nil {
				break
			}
			timestamps = append(timestamps, ts)
		}

		for _, ts := range timestamps {
			lines = append(lines, Line{
				TimeSec: ts,
				Text:    rest,
			})
		}
	}

	for i := range lines {
		lines[i].TimeSec += offsetMs / 1000.0
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i].TimeSec < lines[j].TimeSec
	})

	return lines
}

func parseTimestamp(tag string) (float64, error) {
	parts := strings.SplitN(tag, ":", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid timestamp: %s", tag)
	}
	minutes, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, err
	}
	secStr := strings.TrimSpace(parts[1])
	seconds, err := strconv.ParseFloat(secStr, 64)
	if err != nil {
		return 0, err
	}
	return minutes*60 + seconds, nil
}
