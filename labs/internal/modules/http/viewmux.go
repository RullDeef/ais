package http

import (
	"anicomend/internal/modules/http/dto"
	"anicomend/internal/modules/http/layout"
	"anicomend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const maxPages = 9

type ViewMux struct {
	logger  *zap.SugaredLogger
	service *service.AnimeService
}

func NewViewMux(logger *zap.SugaredLogger, animeService *service.AnimeService) *ViewMux {
	return &ViewMux{
		logger:  logger,
		service: animeService,
	}
}

func (am *ViewMux) AssignHandlers(group *gin.RouterGroup) {
	group.GET("/animes", am.getAnimesPage)
	group.GET("/preferences", am.getPrefernces)
	group.GET("/recomendations", am.getRecomendations)
	group.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/animes")
	})
}

func (am *ViewMux) getAnimesPage(c *gin.Context) {
	currentPage, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		currentPage = 1
	}

	searchQuery := c.Query("query")
	if searchQuery != "" {
		am.service.ApplySearch(searchQuery)
	} else {
		am.service.ClearSearch()
	}

	totalPages := am.service.GetTotalPages()
	animes := am.service.GetPage(int(currentPage))

	animeDTOs := make([]dto.AnimeDTO, len(animes))
	for i, anime := range animes {
		animeDTOs[i] = dto.NewAnimeDTO(anime, am.service.GetPreference(anime.Id))
	}

	c.Status(http.StatusOK)
	err = layout.HomeLayout(c.Writer, layout.HomeLayoutParams{
		Animes:         animeDTOs,
		Pages:          layout.FormatPages(totalPages, maxPages, int(currentPage)),
		FirstPage:      1,
		LastPage:       totalPages,
		SearchQuery:    searchQuery,
		IsSearchResult: searchQuery != "",
		FilterParams:   layout.NewFilterParams(),
	})
	if err != nil {
		panic(err)
	}
}

func (am *ViewMux) getPrefernces(c *gin.Context) {
	animes := am.service.GetPreferencedAnimes()

	animeDTOs := make([]dto.AnimeDTO, len(animes))
	for i, anime := range animes {
		animeDTOs[i] = dto.NewAnimeDTO(anime, am.service.GetPreference(anime.Id))
	}

	c.Status(http.StatusOK)
	layout.PreferencesLayout(c.Writer, layout.PreferencesLayoutParams{
		Animes:        animeDTOs,
		NoPreferences: len(animes) == 0,
	})
}

func (am *ViewMux) getRecomendations(c *gin.Context) {
	currentPage, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		currentPage = 1
	}

	totalPages := am.service.GetRecomendationTotalPages()
	animes := am.service.GetRecomendationPage(int(currentPage))

	animeDTOs := make([]dto.AnimeDTO, len(animes))
	for i, anime := range animes {
		animeDTOs[i] = dto.NewAnimeDTO(anime, am.service.GetPreference(anime.Id))
	}

	c.Status(http.StatusOK)
	layout.RecomendationsLayout(c.Writer, layout.RecomendationsLayoutParams{
		Animes:        animeDTOs,
		Pages:         layout.FormatPages(totalPages, maxPages, int(currentPage)),
		FirstPage:     1,
		LastPage:      totalPages,
		NoPreferences: len(animes) == 0,
		FilterParams:  layout.NewFilterParams(),
	})
}
