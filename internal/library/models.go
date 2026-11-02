package library

type Track struct {
	PersistentID string  `json:"persistentID"`
	Name         string  `json:"name"`
	Artist       string  `json:"artist"`
	Album        string  `json:"album"`
	Duration     float64 `json:"duration"`
	TrackNumber  int     `json:"trackNumber"`
	DiscNumber   int     `json:"discNumber"`
	Genre        string  `json:"genre"`
	Year         int     `json:"year"`
}

type Playlist struct {
	PersistentID string `json:"persistentID"`
	Name         string `json:"name"`
	TrackCount   int    `json:"trackCount"`
	Kind         string `json:"kind"`
	SpecialKind  string `json:"specialKind,omitempty"`
}
