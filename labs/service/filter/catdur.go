package filter

import (
	"anicomend/model"
	"fmt"
)

type DurationCategory int

const (
	VeryShortDuration = DurationCategory(iota)
	ShortDuration
	NotLongDuration
	NotVeryLongDuration
	NotVeryShortDuration
	NotShortDuration
	LongDuration
	VeryLongDuration
)

type durationRange struct {
	min int
	max int
}

var ranges = map[DurationCategory]durationRange{
	VeryShortDuration:    {min: 0, max: 10},
	ShortDuration:        {min: 10, max: 35},
	NotLongDuration:      {min: 35, max: 55},
	NotVeryLongDuration:  {min: 55, max: 81},
	NotVeryShortDuration: {min: 81, max: 96},
	NotShortDuration:     {min: 96, max: 105},
	LongDuration:         {min: 105, max: 130},
	VeryLongDuration:     {min: 130, max: 160},
}

var translations = map[string]DurationCategory{
	"очень короткая":    VeryShortDuration,
	"короткая":          ShortDuration,
	"не долгая":         NotLongDuration,
	"не очень долгая":   NotVeryLongDuration,
	"не очень короткая": NotVeryShortDuration,
	"не короткая":       NotShortDuration,
	"долгая":            LongDuration,
	"очень долгая":      VeryLongDuration,
}

type CategorialDuration struct {
	categories []DurationCategory
}

func NewCategorialDuration() *CategorialDuration {
	return &CategorialDuration{}
}

func (f *CategorialDuration) AddCategory(text string) error {
	if id, ok := translations[text]; ok {
		f.categories = append(f.categories, id)
		return nil
	}
	return fmt.Errorf("no such category '%s'", text)
}

func (f *CategorialDuration) ResetState() {
	f.categories = nil
}

func (f *CategorialDuration) Apply(anime model.Anime) bool {
	for _, id := range f.categories {
		min := float32(ranges[id].min)
		max := float32(ranges[id].max)
		if min <= anime.Duration && anime.Duration <= max {
			return true
		}
	}
	return len(f.categories) == 0
}
