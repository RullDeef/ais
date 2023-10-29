package service

import (
	"anicomend/model"
	"slices"
	"strings"
)

type SearchState struct {
	query   string
	results []model.Anime
}

func NewSearch(query string, animes []model.Anime) *SearchState {
	results := slices.Clone(animes)

	weights := make(map[uint64]float64, len(results))
	for _, a := range results {
		weights[a.Id] = matchWeight(&a, query)
	}

	// sort by match weight
	slices.SortStableFunc(results, func(a, b model.Anime) int {
		a_w := -weights[a.Id]
		b_w := -weights[b.Id]
		if a_w < b_w {
			return -1
		} else if a_w > b_w {
			return 1
		} else {
			return 0
		}
	})

	return &SearchState{
		query:   query,
		results: results,
	}
}

// func matches(anime *model.Anime, query string) bool {
// 	return strings.Contains(anime.Title, query)
// }

// greater is better
func matchWeight(anime *model.Anime, query string) float64 {
	query = strings.ToLower(query)
	title := strings.ToLower(anime.Title)

	bonus := 0
	if strings.Contains(title, query) {
		bonus += len(query)
	}

	return float64(bonus - levenstainDistance(title, query))
}

func levenstainDistance(s1, s2 string) int {
	N1, N2 := len(s1), len(s2)
	if N1 == 0 {
		return N2
	}
	if N2 == 0 {
		return N1
	}
	row1, row2 := make([]int, N1+1), make([]int, N1+1)
	for i := 0; i <= N1; i++ {
		row1[i] = i
	}
	for i := 0; i < N2; i++ {
		row2[0] = i + 1
		for j := 0; j < N1; j++ {
			dis_ins := row2[j] + 1
			dis_rem := row1[j+1] + 1
			dis_rep := row1[j] + 1
			if s1[j] == s2[i] {
				dis_rep--
			}
			row2[j+1] = min(dis_ins, dis_rem, dis_rep)
		}
		row1, row2 = row2, row1
	}
	return row1[N1]
}
