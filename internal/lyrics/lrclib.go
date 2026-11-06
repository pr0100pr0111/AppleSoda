package lyrics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Result struct {
	Lines  []Line
	Plain  string
	Synced bool
}

func Fetch(ctx context.Context, artist, track, album string, durationSec float64) (Result, error) {
	params := url.Values{}
	params.Set("artist_name", artist)
	params.Set("track_name", track)
	params.Set("album_name", album)
	params.Set("duration", strconv.Itoa(int(durationSec)))

	reqURL := "https://lrclib.net/api/get?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return Result{}, fmt.Errorf("create lyrics request: %w", err)
	}
	req.Header.Set("User-Agent", "AppleSoda/0.1 (github.com/pr0100pr0111/AppleSoda)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("fetch lyrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Result{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("lyrics api status %d", resp.StatusCode)
	}

	var body struct {
		SyncedLyrics string `json:"syncedLyrics"`
		PlainLyrics  string `json:"plainLyrics"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return Result{}, fmt.Errorf("decode lyrics response: %w", err)
	}

	if body.SyncedLyrics == "" && body.PlainLyrics == "" {
		return Result{}, nil
	}

	if body.SyncedLyrics != "" {
		parsed := ParseLRC(body.SyncedLyrics)
		if len(parsed) > 0 {
			return Result{Lines: parsed, Synced: true}, nil
		}
	}

	return Result{Plain: body.PlainLyrics}, nil
}
