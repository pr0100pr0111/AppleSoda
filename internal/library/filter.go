package library

import (
	"strings"
)

func FilterTracks(tracks []Track, query string) []Track {
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		res := make([]Track, len(tracks))
		copy(res, tracks)
		return res
	}

	terms := strings.Fields(q)
	var matches []Track

	for _, t := range tracks {
		if matchesAllTerms(t, terms) {
			matches = append(matches, t)
		}
	}
	return matches
}

func FilterIndices(tracks []Track, query string) []int {
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		return nil
	}

	terms := strings.Fields(q)
	var indices []int

	for i, t := range tracks {
		if matchesAllTerms(t, terms) {
			indices = append(indices, i)
		}
	}
	return indices
}

func matchesAllTerms(t Track, terms []string) bool {
	name := strings.ToLower(t.Name)
	artist := strings.ToLower(t.Artist)
	album := strings.ToLower(t.Album)

	for _, term := range terms {
		if !strings.Contains(name, term) && !strings.Contains(artist, term) && !strings.Contains(album, term) {
			return false
		}
	}
	return true
}
