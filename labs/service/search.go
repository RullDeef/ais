package service

import (
	"anicomend/model"
	"strings"
)

type SearchState struct {
	query   string
	results []model.Anime
}

func NewSearch(query string, animes []model.Anime) *SearchState {
	var results []model.Anime
	for _, anime := range animes {
		if matches(&anime, query) {
			results = append(results, anime)
		}
	}
	return &SearchState{
		query:   query,
		results: results,
	}
}

func matches(anime *model.Anime, query string) bool {
	return strings.Contains(anime.Title, query)
}
