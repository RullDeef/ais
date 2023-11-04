package filter

import "anicomend/model"

const (
	minDuration = 0
	maxDuration = 163
)

// duration in minutes
type DurationRangeFilter struct {
	minDuration int
	maxDuration int
}

func NewDurationRangeFilter() *DurationRangeFilter {
	return &DurationRangeFilter{
		minDuration: minDuration,
		maxDuration: maxDuration,
	}
}

func (f *DurationRangeFilter) SetMin(newMin int) {
	if f.maxDuration < newMin {
		f.maxDuration = newMin
	}
	f.minDuration = newMin
}

func (f *DurationRangeFilter) SetMax(newMax int) {
	if f.minDuration > newMax {
		f.minDuration = newMax
	}
	f.maxDuration = newMax
}

func (f *DurationRangeFilter) ResetState() {
	f.minDuration = minDuration
	f.maxDuration = maxDuration
}

func (f *DurationRangeFilter) Apply(anime model.Anime) bool {
	return float32(f.minDuration) <= anime.Duration && anime.Duration <= float32(f.maxDuration)
}
