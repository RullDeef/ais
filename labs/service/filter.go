package service

import (
	"anicomend/model"
	"errors"
)

type FilterMode int

const (
	FilterModeIntersect = FilterMode(iota)
	FilterModeUnion
)

var ErrUnknownFilterMode = errors.New("unknown filter mode")

type FilterBase interface {
	Apply(anime model.Anime) bool
	ResetState()
}

type FilterManager struct {
	bases map[*FilterBase]FilterBase
	mode  FilterMode
}

func NewFilterManager() *FilterManager {
	return &FilterManager{
		bases: make(map[*FilterBase]FilterBase),
		mode:  FilterModeIntersect,
	}
}

func (f *FilterManager) AddFilter(filter FilterBase) {
	f.bases[&filter] = filter
}

func (f *FilterManager) ResetFilters() {
	for _, filter := range f.bases {
		filter.ResetState()
	}
}

func (f *FilterManager) SetFilterMode(mode FilterMode) {
	if mode != FilterModeIntersect && mode != FilterModeUnion {
		panic(ErrUnknownFilterMode)
	}
	f.mode = mode
}

func (f *FilterManager) ApplyFilters(anime model.Anime) bool {
	if f.mode == FilterModeIntersect {
		for _, base := range f.bases {
			if !base.Apply(anime) {
				return false
			}
		}
		return true
	} else if f.mode == FilterModeUnion {
		for _, base := range f.bases {
			if base.Apply(anime) {
				return true
			}
		}
		return len(f.bases) == 0
	} else {
		panic(ErrUnknownFilterMode)
	}
}

func (f *FilterManager) ApplyAllFilters(animes []model.Anime) []model.Anime {
	res := make([]model.Anime, 0, len(animes))
	for _, anime := range animes {
		if f.ApplyFilters(anime) {
			res = append(res, anime)
		}
	}
	return res
}
