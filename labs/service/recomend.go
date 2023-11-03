package service

import (
	"anicomend/model"
	"math"
	"slices"
)

var gs *model.GenreSimilaritier

func recomend(all []model.Anime, prefs []model.PreferenceMark) []model.Anime {
	if gs == nil {
		gs = model.NewGenreSimilaritier()
	}

	if len(prefs) == 0 {
		return nil
	}

	animeMap := make(map[uint64]*model.Anime, len(all))
	for _, a := range all {
		a := a
		animeMap[a.Id] = &a
	}

	// make flow from each preference and find best matches to it
	flows := make([][]model.Anime, len(prefs))
	for i, p := range prefs {
		if p.MarkWeight == model.PreferenceMarkFavourite {
			flows[i] = orderSimilars(animeMap[p.AnimeId], all)
		}
	}

	res := make([]model.Anime, 0)
	resAnimeIds := make(map[uint64]struct{}, len(all))

	// do not put animes from prefs in result
	for _, p := range prefs {
		resAnimeIds[p.AnimeId] = struct{}{}
	}

	// interleave flows (simple way)
	for i := 0; i < len(all); i++ {
		for f := 0; f < len(flows); f++ {
			if prefs[f].MarkWeight == model.PreferenceMarkUnfavourite {
				continue
			}
			a := flows[f][i]
			if _, ok := resAnimeIds[a.Id]; !ok {
				res = append(res, a)
				resAnimeIds[a.Id] = struct{}{}
			}
		}
	}

	return res
}

func orderSimilars(base *model.Anime, all []model.Anime) []model.Anime {
	all = slices.Clone(all)
	slices.SortStableFunc(all, func(a1, a2 model.Anime) int {
		m1 := similarityMetric(*base, a1)
		m2 := similarityMetric(*base, a2)
		if m1 < m2 {
			return 1
		} else if m1 > m2 {
			return -1
		} else {
			return 0
		}
	})
	return all
}

func similarityMetric(a1, a2 model.Anime) float64 {
	m1 := similarityMetric_1(a1, a2)
	m2 := similarityMetric_2(a1, a2)
	m3 := similarityMetric_3(a1, a2)

	w1, w2, w3 := 0.9, 0.5, 0.2
	w123 := w1 + w2 + w3

	return (0.9*m1 + 0.5*m2 + 0.2*m3) / w123
}

// basic metric 1 (by genres)
func similarityMetric_1(a1, a2 model.Anime) float64 {
	return gs.Similarity(a1.Genres, a2.Genres)
}

// basic metric 2 (by aired dates)
func similarityMetric_2(a1, a2 model.Anime) float64 {
	t := a2.AiredFrom.Sub(a1.AiredFrom).Hours() / 24
	return math.Exp(-math.Abs(t) / 900)
}

// basic metric 3 (by episodes count)
func similarityMetric_3(a1, a2 model.Anime) float64 {
	t := math.Abs(float64(a2.Episodes - a1.Episodes))
	return math.Exp(-t / 10.0)
}
