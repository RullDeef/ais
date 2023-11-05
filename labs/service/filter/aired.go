package filter

import (
	"anicomend/model"
	"time"
)

var (
	MinAiredTime = time.Date(1942, 03, 15, 0, 0, 0, 0, time.UTC)
	MaxAiredTime = time.Date(2018, 05, 25, 0, 0, 0, 0, time.UTC)
)

type AiredRangeFilter struct {
	minTime time.Time
	maxTime time.Time
}

func NewAiredRangeFilter() *AiredRangeFilter {
	return &AiredRangeFilter{
		minTime: MinAiredTime,
		maxTime: MaxAiredTime,
	}
}

func (f *AiredRangeFilter) SetMin(newMin time.Time) {
	if f.maxTime.Before(newMin) {
		f.maxTime = newMin
	}
	f.minTime = newMin
}

func (f *AiredRangeFilter) SetMax(newMax time.Time) {
	if f.minTime.After(newMax) {
		f.minTime = newMax
	}
	f.maxTime = newMax
}

func (f *AiredRangeFilter) SetMinYear(year int) {
	f.SetMin(time.Date(year, 0, 0, 0, 0, 0, 0, time.UTC))
}

func (f *AiredRangeFilter) SetMaxYear(year int) {
	f.SetMax(time.Date(year, 0, 0, 0, 0, 0, 0, time.UTC))
}

func (f *AiredRangeFilter) Apply(anime model.Anime) bool {
	return !f.minTime.After(anime.AiredFrom) && !anime.AiredFrom.After(f.maxTime)
}

func (f *AiredRangeFilter) ResetState() {
	f.minTime = MinAiredTime
	f.maxTime = MaxAiredTime
}
