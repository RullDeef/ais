package filter

import (
	"anicomend/model"
	"slices"
)

type SimpleGenreFilter struct {
	selected []string
}

func NewSimpleGenreFilter() *SimpleGenreFilter {
	return &SimpleGenreFilter{}
}

func (f *SimpleGenreFilter) Select(genre string) {
	if !slices.Contains(f.selected, genre) {
		f.selected = append(f.selected, genre)
	}
}

func (f *SimpleGenreFilter) Deselect(genre string) {
	if i := slices.Index(f.selected, genre); i != -1 {
		f.selected = slices.Delete(f.selected, i, i+1)
	}
}

func (f *SimpleGenreFilter) Apply(anime model.Anime) bool {
	for _, genre := range f.selected {
		if !slices.Contains(anime.Genres, genre) {
			return false
		}
	}
	return true
}

func (f *SimpleGenreFilter) ResetState() {
	f.selected = nil
}
