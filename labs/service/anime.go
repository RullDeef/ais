package service

import (
	"anicomend/model"
	"strings"
)

const ItemsPerPage = 12

type AnimeService struct {
	animes []model.Anime
}

func NewAnimeService(loader model.AnimeLoader) (*AnimeService, error) {
	animes, err := loader.Load()
	if err != nil {
		return nil, err
	}

	// filter adult anime
	animes = filterAdultGenres(animes)

	// transform image cdn
	for i, anime := range animes {
		animes[i].ImageURL = strings.ReplaceAll(anime.ImageURL, "myanimelist.cdn-dena.com", "cdn.myanimelist.net")
	}

	return &AnimeService{
		animes: animes,
	}, nil
}

func (a *AnimeService) GetTotalPages() int {
	return (len(a.animes) + ItemsPerPage - 1) / ItemsPerPage
}

func (a *AnimeService) GetPage(page int) []model.Anime {
	return a.animes[(page-1)*ItemsPerPage : page*ItemsPerPage]
}

func filterAdultGenres(animes []model.Anime) []model.Anime {
	filtered := make([]model.Anime, 0)
	for _, anime := range animes {
		hasAdultGenre := false
		for _, genre := range anime.Genres {
			if genreIsAdult(genre) {
				hasAdultGenre = true
				break
			}
		}
		if !hasAdultGenre {
			filtered = append(filtered, anime)
		}
	}
	return filtered
}

func genreIsAdult(genre string) bool {
	return genre == "Hentai" || genre == "Yaoi" || genre == "Yuri"
}
