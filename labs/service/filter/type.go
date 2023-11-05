package filter

import (
	"anicomend/model"
	"errors"
	"slices"
)

var (
	ErrUnknownAnimeType = errors.New("unknown anime type")
)

var AnimeTypes = []string{
	"TV",
	"OVA",
	"Special",
	"ONA",
	"Movie",
	"Music",
}

type SimpleTypeFilter struct {
	selected []string
}

func NewSimpleTypeFilter() *SimpleTypeFilter {
	return &SimpleTypeFilter{}
}

func (f *SimpleTypeFilter) Select(animeType string) {
	if !slices.Contains(AnimeTypes, animeType) {
		panic(ErrUnknownAnimeType)
	}
	if !slices.Contains(f.selected, animeType) {
		f.selected = append(f.selected, animeType)
	}
}

func (f *SimpleTypeFilter) ResetState() {
	f.selected = nil
}

func (f *SimpleTypeFilter) Apply(anime model.Anime) bool {
	return len(f.selected) == 0 || slices.Contains(f.selected, anime.Type)
}
