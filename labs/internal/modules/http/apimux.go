package http

import (
	"anicomend/service"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApiMux struct {
	logger  *zap.SugaredLogger
	service *service.AnimeService
}

func NewApiMux(logger *zap.SugaredLogger, service *service.AnimeService) *ApiMux {
	return &ApiMux{
		logger:  logger,
		service: service,
	}
}

func (am *ApiMux) AssignHandlers(routerGroup *gin.RouterGroup) {
	// /api/animes/{anime_id}?mark=[fav,unfav,clear]
	routerGroup.GET("/animes/:anime_id", am.updateMark)
	routerGroup.GET("/animes/:anime_id/", am.updateMark)

	// /api/animes?page=N
	routerGroup.GET("/animes", am.getPage)

	// /api/filter POST method
	routerGroup.POST("/filter", am.applyFilters)
}

func (am *ApiMux) getPage(c *gin.Context) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil || page < 1 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	animes := am.service.GetPage(int(page))
	c.Status(http.StatusOK)
	json.NewEncoder(c.Writer).Encode(animes)
}

func (am *ApiMux) updateMark(c *gin.Context) {
	animeId, err := strconv.ParseUint(c.Param("anime_id"), 10, 64)
	if err != nil {
		am.logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	switch mark := c.Query("mark"); mark {
	case "fav":
		am.service.MarkAsFavorite(animeId)
		c.Status(http.StatusNoContent)
	case "unfav":
		am.service.MarkAsUnfavorite(animeId)
		c.Status(http.StatusNoContent)
	case "clear":
		am.service.ClearPreferenceMark(animeId)
	default:
		am.logger.Error("failed to parse mark:", mark)
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (am *ApiMux) applyFilters(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	println(string(body))
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	var form FiltersForm
	c.Bind(&form)

	fmt.Printf("%+v\n", form)

	am.service.GenreFilter.ResetState()
	for _, genre := range form.Genres {
		am.service.GenreFilter.Select(genre)
	}

	am.service.DurationFilter.ResetState()
	am.service.DurationFilter.SetMin(form.DurationMin)
	am.service.DurationFilter.SetMax(form.DurationMax)

	am.service.CatDurFilter.ResetState()
	for _, cat := range form.DurationCat {
		if err := am.service.CatDurFilter.AddCategory(cat); err != nil {
			am.logger.Error(err)
		}
	}

	am.service.AiredFilter.ResetState()
	am.service.AiredFilter.SetMinYear(form.AiredMin)
	am.service.AiredFilter.SetMaxYear(form.AiredMax)

	am.service.TypeFilter.ResetState()
	for _, animeType := range form.Types {
		am.service.TypeFilter.Select(animeType)
	}

	c.Status(http.StatusOK)
}

type FiltersForm struct {
	Genres      []string `form:"genre"`
	DurationMin int      `form:"duration-min"`
	DurationMax int      `form:"duration-max"`
	DurationCat []string `form:"duration-cat"`
	AiredMin    int      `form:"aired-min"`
	AiredMax    int      `form:"aired-max"`
	Types       []string `form:"type"`
}
