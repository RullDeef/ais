package model

const (
	PreferenceMarkFavourite   = 1
	PreferenceMarkUnfavourite = -1
)

type Anime struct {
	Id       uint64
	Title    string
	ImageURL string
	Type     string
	Source   string
	Studio   string
	Genres   []string
	Duration float32
	Year     uint32
}

type PreferenceMark struct {
	AnimeId    uint64
	MarkWeight int
}

type AnimeLoader interface {
	Load() ([]Anime, error)
}
