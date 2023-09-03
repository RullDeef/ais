package model

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

type AnimeLoader interface {
	Load() ([]Anime, error)
}
