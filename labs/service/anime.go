package service

import (
	"anicomend/model"
	"slices"
	"strings"
)

const ItemsPerPage = 12

type AnimeService struct {
	animes          []model.Anime
	preferenceMarks []model.PreferenceMark
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

func (a *AnimeService) MarkAsFavorite(animeId uint64) {
	// check if mark already set
	for i, anime := range a.preferenceMarks {
		if anime.AnimeId == animeId {
			a.preferenceMarks[i].MarkWeight = model.PreferenceMarkFavourite
			return
		}
	}
	for _, anime := range a.animes {
		if anime.Id == animeId {
			a.preferenceMarks = append(a.preferenceMarks, model.PreferenceMark{
				AnimeId:    animeId,
				MarkWeight: model.PreferenceMarkFavourite,
			})
			break
		}
	}
}

func (a *AnimeService) MarkAsUnfavorite(animeId uint64) {
	// check if mark already set
	for i, anime := range a.preferenceMarks {
		if anime.AnimeId == animeId {
			a.preferenceMarks[i].MarkWeight = model.PreferenceMarkUnfavourite
			return
		}
	}
	for _, anime := range a.animes {
		if anime.Id == animeId {
			a.preferenceMarks = append(a.preferenceMarks, model.PreferenceMark{
				AnimeId:    animeId,
				MarkWeight: model.PreferenceMarkUnfavourite,
			})
			break
		}
	}
}

func (a *AnimeService) ClearPreferenceMark(animeId uint64) {
	for i, anime := range a.preferenceMarks {
		if anime.AnimeId == animeId {
			a.preferenceMarks = slices.Delete(a.preferenceMarks, i, i+1)
			break
		}
	}
}

func (a *AnimeService) ClearAllPreferences() {
	a.preferenceMarks = nil
}

func (a *AnimeService) GetTotalPages() int {
	return (len(a.animes) + ItemsPerPage - 1) / ItemsPerPage
}

func (a *AnimeService) GetPage(page int) []model.Anime {
	return a.animes[(page-1)*ItemsPerPage : page*ItemsPerPage]
}

func (a *AnimeService) GetPreference(animeId uint64) model.PreferenceMark {
	for _, pref := range a.preferenceMarks {
		if pref.AnimeId == animeId {
			return pref
		}
	}
	return model.PreferenceMark{
		AnimeId:    animeId,
		MarkWeight: 0,
	}
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
