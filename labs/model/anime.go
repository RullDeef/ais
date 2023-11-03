package model

import "time"

const (
	PreferenceMarkFavourite   = 1
	PreferenceMarkUnfavourite = -1
)

type Anime struct {
	Id        uint64
	Title     string
	ImageURL  string
	Type      string
	Source    string
	Studio    string
	Genres    []string
	Duration  float32
	Episodes  uint32
	Year      uint32
	AiredFrom time.Time
	AiredTo   time.Time
}

type PreferenceMark struct {
	AnimeId    uint64
	MarkWeight int
}

type AnimeLoader interface {
	Load() ([]Anime, error)
}
