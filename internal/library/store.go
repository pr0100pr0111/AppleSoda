package library

import (
	"sync"
)

type Store struct {
	mu        sync.RWMutex
	cacheDir  string
	playlists []Playlist
	tracks    map[string][]Track
}

func NewStore(cacheDir string) *Store {
	return &Store{
		cacheDir: cacheDir,
		tracks:   make(map[string][]Track),
	}
}

func (s *Store) LoadFromDisk() error {
	pls, trs, err := LoadCache(s.cacheDir)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(pls) > 0 {
		s.playlists = pls
	}
	if trs != nil {
		s.tracks = trs
	}
	return nil
}

func (s *Store) SaveToDisk() error {
	s.mu.RLock()
	pls := make([]Playlist, len(s.playlists))
	copy(pls, s.playlists)
	trs := make(map[string][]Track, len(s.tracks))
	for k, v := range s.tracks {
		copied := make([]Track, len(v))
		copy(copied, v)
		trs[k] = copied
	}
	s.mu.RUnlock()

	return SaveCache(s.cacheDir, pls, trs)
}

func (s *Store) GetPlaylists() []Playlist {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]Playlist, len(s.playlists))
	copy(res, s.playlists)
	return res
}

func (s *Store) SetPlaylists(playlists []Playlist) {
	s.mu.Lock()
	s.playlists = make([]Playlist, len(playlists))
	copy(s.playlists, playlists)
	s.mu.Unlock()

	_ = s.SaveToDisk()
}

func (s *Store) GetTracks(playlistID string) ([]Track, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tracks[playlistID]
	if !ok {
		return nil, false
	}
	res := make([]Track, len(t))
	copy(res, t)
	return res, true
}

func (s *Store) SetTracks(playlistID string, tracks []Track) {
	s.mu.Lock()
	copied := make([]Track, len(tracks))
	copy(copied, tracks)
	s.tracks[playlistID] = copied
	s.mu.Unlock()

	_ = s.SaveToDisk()
}

func (s *Store) HasTracks(playlistID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.tracks[playlistID]
	return ok
}

func (s *Store) TrackCount(playlistID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if t, ok := s.tracks[playlistID]; ok {
		return len(t)
	}
	return 0
}
