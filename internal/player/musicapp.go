package player

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type MusicAppPlayer struct {
	autoLaunch bool
}

func NewMusicAppPlayer(autoLaunch bool) *MusicAppPlayer {
	return &MusicAppPlayer{
		autoLaunch: autoLaunch,
	}
}

func (m *MusicAppPlayer) IsRunning(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "pgrep", "-x", "Music")
	return cmd.Run() == nil
}

func (m *MusicAppPlayer) Launch(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "open", "-ga", "Music")
	return cmd.Run()
}

func (m *MusicAppPlayer) checkRunning(ctx context.Context) error {
	if !m.IsRunning(ctx) {
		if m.autoLaunch {
			if err := m.Launch(ctx); err != nil {
				return fmt.Errorf("launch music.app: %w", err)
			}
			select {
			case <-time.After(1 * time.Second):
			case <-ctx.Done():
				return ctx.Err()
			}
		} else {
			return ErrMusicAppNotRunning
		}
	}
	return nil
}

type jxaStatusResponse struct {
	State    string  `json:"state"`
	Position float64 `json:"position"`
	Volume   int     `json:"volume"`
	Shuffle  bool    `json:"shuffle"`
	Repeat   string  `json:"repeat"`
	Track    *struct {
		Name         string  `json:"name"`
		Artist       string  `json:"artist"`
		Album        string  `json:"album"`
		Duration     float64 `json:"duration"`
		PersistentID string  `json:"persistentID"`
	} `json:"track"`
	Playlist *struct {
		Name         string `json:"name"`
		PersistentID string `json:"persistentID"`
	} `json:"playlist"`
	Error string `json:"error"`
}

func (m *MusicAppPlayer) Status(ctx context.Context) (Status, error) {
	if err := m.checkRunning(ctx); err != nil {
		return Status{}, err
	}

	script := `(() => {
		try {
			const app = Application("Music");
			const state = app.playerState();
			let pos = 0;
			try { pos = app.playerPosition(); } catch(e) {}
			const vol = app.soundVolume();
			const shuf = app.shuffleEnabled();
			const rep = app.songRepeat();
			let trackInfo = null;
			try {
				const tr = app.currentTrack();
				if (tr && tr.name()) {
					trackInfo = {
						name: tr.name(),
						artist: tr.artist() || "",
						album: tr.album() || "",
						duration: tr.duration() || 0,
						persistentID: tr.persistentID() || ""
					};
				}
			} catch(e) {}
			let plInfo = null;
			try {
				const pl = app.currentPlaylist();
				if (pl && pl.name()) {
					plInfo = { name: pl.name(), persistentID: pl.persistentID() };
				}
			} catch(e) {}
			return JSON.stringify({
				state: state,
				position: pos,
				volume: vol,
				shuffle: shuf,
				repeat: rep,
				track: trackInfo,
				playlist: plInfo
			});
		} catch(err) {
			return JSON.stringify({ error: err.message });
		}
	})()`

	out, err := runJXA(ctx, script)
	if err != nil {
		return Status{}, err
	}

	var res jxaStatusResponse
	if err := json.Unmarshal(out, &res); err != nil {
		return Status{}, fmt.Errorf("parse status response: %w", err)
	}
	if res.Error != "" {
		return Status{}, fmt.Errorf("music status: %s", res.Error)
	}

	st := Status{
		State:    ParsePlayerState(res.State),
		Position: res.Position,
		Volume:   res.Volume,
		Shuffle:  res.Shuffle,
		Repeat:   ParseRepeatMode(res.Repeat),
	}

	if res.Track != nil && res.Track.PersistentID != "" {
		st.HasTrack = true
		st.Track = Track{
			PersistentID: res.Track.PersistentID,
			Name:         res.Track.Name,
			Artist:       res.Track.Artist,
			Album:        res.Track.Album,
			Duration:     res.Track.Duration,
		}
	}

	if res.Playlist != nil {
		st.PlaylistID = res.Playlist.PersistentID
		st.PlaylistName = res.Playlist.Name
	}

	return st, nil
}

func (m *MusicAppPlayer) Play(ctx context.Context) error {
	if err := m.checkRunning(ctx); err != nil {
		return err
	}
	_, err := runJXA(ctx, `Application("Music").play();`)
	return err
}

func (m *MusicAppPlayer) Pause(ctx context.Context) error {
	if err := m.checkRunning(ctx); err != nil {
		return err
	}
	_, err := runJXA(ctx, `Application("Music").pause();`)
	return err
}

func (m *MusicAppPlayer) TogglePlayPause(ctx context.Context) error {
	if err := m.checkRunning(ctx); err != nil {
		return err
	}
	_, err := runJXA(ctx, `Application("Music").playpause();`)
	return err
}

func (m *MusicAppPlayer) Next(ctx context.Context) error {
	if err := m.checkRunning(ctx); err != nil {
		return err
	}
	_, err := runJXA(ctx, `Application("Music").nextTrack();`)
	return err
}

func (m *MusicAppPlayer) Previous(ctx context.Context) error {
	if err := m.checkRunning(ctx); err != nil {
		return err
	}
	_, err := runJXA(ctx, `Application("Music").previousTrack();`)
	return err
}

func (m *MusicAppPlayer) Seek(ctx context.Context, position float64) error {
	if err := m.checkRunning(ctx); err != nil {
		return err
	}
	if position < 0 {
		position = 0
	}
	script := fmt.Sprintf(`Application("Music").playerPosition = %f;`, position)
	_, err := runJXA(ctx, script)
	return err
}

func (m *MusicAppPlayer) SetVolume(ctx context.Context, volume int) error {
	if err := m.checkRunning(ctx); err != nil {
		return err
	}
	if volume < 0 {
		volume = 0
	}
	if volume > 100 {
		volume = 100
	}
	script := fmt.Sprintf(`Application("Music").soundVolume = %d;`, volume)
	_, err := runJXA(ctx, script)
	return err
}

func (m *MusicAppPlayer) SetShuffle(ctx context.Context, on bool) error {
	if err := m.checkRunning(ctx); err != nil {
		return err
	}
	script := fmt.Sprintf(`Application("Music").shuffleEnabled = %t;`, on)
	_, err := runJXA(ctx, script)
	return err
}

func (m *MusicAppPlayer) SetRepeat(ctx context.Context, mode RepeatMode) error {
	if err := m.checkRunning(ctx); err != nil {
		return err
	}
	modeStr := "off"
	switch mode {
	case RepeatAll:
		modeStr = "all"
	case RepeatOne:
		modeStr = "one"
	}
	script := fmt.Sprintf(`Application("Music").songRepeat = "%s";`, modeStr)
	_, err := runJXA(ctx, script)
	return err
}

type jxaGenericResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (m *MusicAppPlayer) PlayTrack(ctx context.Context, playlistID string, trackID string) error {
	if err := m.checkRunning(ctx); err != nil {
		return err
	}
	script := fmt.Sprintf(`(() => {
		try {
			const app = Application("Music");
			let tr = null;
			if ("%s" !== "") {
				const pls = app.playlists.whose({ persistentID: "%s" });
				if (pls.length > 0) {
					const trs = pls[0].tracks.whose({ persistentID: "%s" });
					if (trs.length > 0) tr = trs[0];
				}
			}
			if (!tr) {
				const trs = app.tracks.whose({ persistentID: "%s" });
				if (trs.length > 0) tr = trs[0];
			}
			if (!tr) return JSON.stringify({ error: "track not found" });
			app.play(tr);
			return JSON.stringify({ success: true });
		} catch(e) {
			return JSON.stringify({ error: e.message });
		}
	})()`, playlistID, playlistID, trackID, trackID)

	out, err := runJXA(ctx, script)
	if err != nil {
		return err
	}

	var res jxaGenericResult
	if err := json.Unmarshal(out, &res); err != nil {
		return fmt.Errorf("parse play track response: %w", err)
	}
	if res.Error != "" {
		return fmt.Errorf("%w: %s", ErrTrackUnplayable, res.Error)
	}
	return nil
}

type jxaPlaylist struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (m *MusicAppPlayer) ListPlaylists(ctx context.Context) ([]Playlist, error) {
	if err := m.checkRunning(ctx); err != nil {
		return nil, err
	}
	script := `(() => {
		try {
			const app = Application("Music");
			const pls = app.playlists;
			const names = pls.name();
			const ids = pls.persistentID();
			const res = [];
			for (let i = 0; i < names.length; i++) {
				if (names[i]) {
					res.push({ name: names[i], id: ids[i] || "" });
				}
			}
			return JSON.stringify(res);
		} catch(e) {
			return JSON.stringify({ error: e.message });
		}
	})()`

	out, err := runJXA(ctx, script)
	if err != nil {
		return nil, err
	}

	var rawList []jxaPlaylist
	if err := json.Unmarshal(out, &rawList); err != nil {
		var errObj jxaGenericResult
		if err2 := json.Unmarshal(out, &errObj); err2 == nil && errObj.Error != "" {
			return nil, fmt.Errorf("list playlists: %s", errObj.Error)
		}
		return nil, fmt.Errorf("parse playlists: %w", err)
	}

	result := make([]Playlist, 0, len(rawList))
	for _, p := range rawList {
		result = append(result, Playlist{
			PersistentID: p.ID,
			Name:         p.Name,
		})
	}
	return result, nil
}

type jxaTracksResponse struct {
	Total  int `json:"total"`
	Tracks []struct {
		Name         string  `json:"name"`
		Artist       string  `json:"artist"`
		Album        string  `json:"album"`
		Duration     float64 `json:"duration"`
		PersistentID string  `json:"persistentID"`
		Index        int     `json:"index"`
	} `json:"tracks"`
	Error string `json:"error"`
}

func (m *MusicAppPlayer) ListTracks(ctx context.Context, playlistID string, offset, limit int) ([]Track, int, error) {
	if err := m.checkRunning(ctx); err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}

	script := fmt.Sprintf(`(() => {
		try {
			const app = Application("Music");
			let container = null;
			if ("%s" !== "") {
				const pls = app.playlists.whose({ persistentID: "%s" });
				if (pls.length > 0) container = pls[0];
			}
			if (!container) {
				container = app.playlists.byName("Медиатека") || app.playlists.byName("Library") || app.playlists[0];
			}
			if (!container) return JSON.stringify({ total: 0, tracks: [] });
			const trs = container.tracks;
			const names = trs.name();
			const total = names.length;
			const artists = trs.artist();
			const albums = trs.album();
			const durs = trs.duration();
			const ids = trs.persistentID();
			const start = %d;
			const limit = %d;
			const end = Math.min(total, start + limit);
			const list = [];
			for (let i = start; i < end; i++) {
				list.push({
					name: names[i] || "",
					artist: artists[i] || "",
					album: albums[i] || "",
					duration: durs[i] || 0,
					persistentID: ids[i] || "",
					index: i + 1
				});
			}
			return JSON.stringify({ total: total, tracks: list });
		} catch(e) {
			return JSON.stringify({ error: e.message });
		}
	})()`, playlistID, playlistID, offset, limit)

	out, err := runJXA(ctx, script)
	if err != nil {
		return nil, 0, err
	}

	var res jxaTracksResponse
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, 0, fmt.Errorf("parse tracks: %w", err)
	}
	if res.Error != "" {
		return nil, 0, fmt.Errorf("list tracks: %s", res.Error)
	}

	tracks := make([]Track, len(res.Tracks))
	for i, t := range res.Tracks {
		tracks[i] = Track{
			PersistentID: t.PersistentID,
			Name:         t.Name,
			Artist:       t.Artist,
			Album:        t.Album,
			Duration:     t.Duration,
			TrackNumber:  t.Index,
		}
	}

	return tracks, res.Total, nil
}

type jxaArtworkResponse struct {
	Hex   string `json:"hex"`
	Error string `json:"error"`
}

func (m *MusicAppPlayer) Artwork(ctx context.Context) ([]byte, string, error) {
	if err := m.checkRunning(ctx); err != nil {
		return nil, "", err
	}

	script := `(() => {
		try {
			const app = Application("Music");
			const tr = app.currentTrack();
			if (!tr) return JSON.stringify({ error: "no track" });
			const arts = tr.artworks;
			if (!arts || arts.length === 0) return JSON.stringify({ error: "no artwork" });
			const raw = arts[0].rawData();
			if (!raw) return JSON.stringify({ error: "empty artwork" });
			const start = raw.indexOf("$");
			if (start === -1) return JSON.stringify({ error: "unknown format" });
			const end = raw.lastIndexOf(")");
			const hex = raw.substring(start + 1, end !== -1 ? end : raw.length);
			return JSON.stringify({ hex: hex });
		} catch(e) {
			return JSON.stringify({ error: e.message });
		}
	})()`

	out, err := runJXA(ctx, script)
	if err != nil {
		return nil, "", err
	}

	var res jxaArtworkResponse
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, "", fmt.Errorf("parse artwork response: %w", err)
	}
	if res.Error != "" || res.Hex == "" {
		return nil, "", ErrArtworkNotFound
	}

	data, err := hex.DecodeString(res.Hex)
	if err != nil {
		return nil, "", fmt.Errorf("decode artwork hex: %w", err)
	}

	format := "image/jpeg"
	if bytes.HasPrefix(data, []byte{0x89, 'P', 'N', 'G'}) {
		format = "image/png"
	}

	return data, format, nil
}
