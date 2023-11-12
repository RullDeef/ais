package service

import (
	"anicomend/model"
	"anicomend/service/filter"
	"errors"
	"slices"

	"go.uber.org/zap"
)

const (
	ItemsPerPage = 12
)

type AnimeService struct {
	*FilterManager

	animes []model.Anime

	// specific filters
	GenreFilter    *filter.SimpleGenreFilter
	DurationFilter *filter.DurationRangeFilter
	CatDurFilter   *filter.CategorialDuration
	AiredFilter    *filter.AiredRangeFilter
	TypeFilter     *filter.SimpleTypeFilter

	preferenceMarks []model.PreferenceMark
	recomended      []model.Anime

	// search works independently from filters and recomendations
	searchState *SearchState

	logger *zap.SugaredLogger
}

func NewAnimeService(loader model.AnimeLoader, logger *zap.SugaredLogger) (*AnimeService, error) {
	animes, err := loader.Load()
	if err != nil {
		return nil, err
	}

	service := AnimeService{
		FilterManager:  NewFilterManager(),
		animes:         animes,
		GenreFilter:    filter.NewSimpleGenreFilter(),
		DurationFilter: filter.NewDurationRangeFilter(),
		CatDurFilter:   filter.NewCategorialDuration(),
		AiredFilter:    filter.NewAiredRangeFilter(),
		TypeFilter:     filter.NewSimpleTypeFilter(),
		logger:         logger,
	}

	service.AddFilter(service.GenreFilter)
	service.AddFilter(service.DurationFilter)
	service.AddFilter(service.CatDurFilter)
	service.AddFilter(service.AiredFilter)
	service.AddFilter(service.TypeFilter)

	return &service, nil
}

func (a *AnimeService) getFilteredAnimes() []model.Anime {
	return a.ApplyAllFilters(a.animes)
}

func (a *AnimeService) GetTotalPages() int {
	if a.searchState != nil {
		return a.GetSearchTotalPages()
	}
	return (len(a.getFilteredAnimes()) + ItemsPerPage - 1) / ItemsPerPage
}

func (a *AnimeService) GetPage(page int) []model.Anime {
	if a.searchState != nil {
		return a.GetSearchPage(page)
	}
	animes := a.getFilteredAnimes()
	page = min(page, a.GetTotalPages())
	lower := max((page-1)*ItemsPerPage, 0)
	upper := min(page*ItemsPerPage, len(animes))
	return animes[lower:upper]
}

// Preferences handling

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
				Anime:      anime,
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
				Anime:      anime,
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

func (a *AnimeService) GetPreference(animeId uint64) model.PreferenceMark {
	for _, pref := range a.preferenceMarks {
		if pref.AnimeId == animeId {
			return pref
		}
	}
	for _, anime := range a.animes {
		if anime.Id == animeId {
			return model.PreferenceMark{
				AnimeId:    animeId,
				Anime:      anime,
				MarkWeight: 0,
			}
		}
	}
	panic(errors.New("failed to find anime by id"))
}

func (a *AnimeService) GetPreferencedAnimes() []model.Anime {
	res := make([]model.Anime, 0, len(a.preferenceMarks))
	for _, mark := range a.preferenceMarks {
		res = append(res, mark.Anime)
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
	a.regenerateRecomendations()
	return (len(a.recomended) + ItemsPerPage - 1) / ItemsPerPage
}

func (a *AnimeService) GetRecomendationPage(page int) []model.Anime {
	a.regenerateRecomendations()
	page = min(page, max(1, a.GetRecomendationTotalPages()))
	lower := max(0, (page-1)*ItemsPerPage)
	upper := min(page*ItemsPerPage, len(a.recomended))
	return a.recomended[lower:upper]
}

func (a *AnimeService) regenerateRecomendations() {
	var err error
	a.recomended, err = recomend(a.getFilteredAnimes(), a.preferenceMarks)
	if err != nil {
		a.logger.Error(err)
	}
}
