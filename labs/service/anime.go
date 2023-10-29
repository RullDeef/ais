package service

import (
	"anicomend/model"
	"slices"
	"strings"
)

const (
	ItemsPerPage = 12
)

var filteredGenres []string = []string{
	"Dementia",
	"Harem",
	"Kids",
	"Ecchi",
	"Shounen Ai",
	"Yuri",
	"Hentai",
	"Yaoi",
	"Shoujo Ai",
}

type AnimeService struct {
	animes          []model.Anime
	preferenceMarks []model.PreferenceMark

	recomended      []model.Anime
	recomendedDirty bool

	searchState *SearchState
}

func NewAnimeService(loader model.AnimeLoader) (*AnimeService, error) {
	animes, err := loader.Load()
	if err != nil {
		return nil, err
	}

	// filter adult anime
	animes = filterGenres(animes)

	// transform image cdn
	for i, anime := range animes {
		animes[i].ImageURL = strings.ReplaceAll(anime.ImageURL, "myanimelist.cdn-dena.com", "cdn.myanimelist.net")
	}

	return &AnimeService{
		animes:          animes,
		recomendedDirty: true,
	}, nil
}

func filterGenres(animes []model.Anime) []model.Anime {
	filtered := make([]model.Anime, 0)
	for _, anime := range animes {
		hasFilteredGenre := false
		for _, genre := range anime.Genres {
			if slices.Contains(filteredGenres, genre) {
				hasFilteredGenre = true
				break
			}
		}
		if !hasFilteredGenre && len(anime.Genres) > 0 {
			filtered = append(filtered, anime)
		}
	}
	return filtered
}

func (a *AnimeService) GetTotalPages() int {
	if a.searchState != nil {
		return a.GetSearchTotalPages()
	}
	return (len(a.animes) + ItemsPerPage - 1) / ItemsPerPage
}

func (a *AnimeService) GetPage(page int) []model.Anime {
	if a.searchState != nil {
		return a.GetSearchPage(page)
	}
	upper := min(page*ItemsPerPage, len(a.animes))
	return a.animes[(page-1)*ItemsPerPage : upper]
}

// Preferences handling

func (a *AnimeService) MarkAsFavorite(animeId uint64) {
	a.recomendedDirty = true
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
	a.recomendedDirty = true
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
	a.recomendedDirty = true
	for i, anime := range a.preferenceMarks {
		if anime.AnimeId == animeId {
			a.preferenceMarks = slices.Delete(a.preferenceMarks, i, i+1)
			break
		}
	}
}

func (a *AnimeService) ClearAllPreferences() {
	a.recomendedDirty = true
	a.preferenceMarks = nil
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

func (a *AnimeService) GetPreferencedAnimes() []model.Anime {
	res := make([]model.Anime, 0, len(a.preferenceMarks))
	for _, mark := range a.preferenceMarks {
		for _, anime := range a.animes {
			if anime.Id == mark.AnimeId {
				res = append(res, anime)
			}
		}
	}
	return res
}

// Search state handling

func (a *AnimeService) ClearSearch() {
	a.searchState = nil
}

func (a *AnimeService) ApplySearch(query string) {
	if a.searchState == nil || a.searchState.query != query {
		a.searchState = NewSearch(query, a.animes)
	}
}

func (a *AnimeService) GetSearchTotalPages() int {
	if a.searchState == nil {
		return 0
	}
	return (len(a.searchState.results) + ItemsPerPage - 1) / ItemsPerPage
}

func (a *AnimeService) GetSearchPage(page int) []model.Anime {
	if a.searchState == nil {
		return nil
	}
	upper := min(page*ItemsPerPage, len(a.searchState.results))
	return a.searchState.results[(page-1)*ItemsPerPage : upper]
}

// recomendations

func (a *AnimeService) GetRecomendationTotalPages() int {
	if a.recomendedDirty {
		a.regenerateRecomendations()
		a.recomendedDirty = false
	}
	return (len(a.recomended) + ItemsPerPage - 1) / ItemsPerPage
}

func (a *AnimeService) GetRecomendationPage(page int) []model.Anime {
	if a.recomendedDirty {
		a.regenerateRecomendations()
		a.recomendedDirty = false
	}
	page = min(page, max(1, a.GetRecomendationTotalPages()))
	upper := min(page*ItemsPerPage, len(a.recomended))
	return a.recomended[(page-1)*ItemsPerPage : upper]
}

func (a *AnimeService) regenerateRecomendations() {
	a.recomended = recomend(a.animes, a.preferenceMarks)
}
