package http

import (
	"anicomend/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AnimeMux struct {
	logger  *zap.SugaredLogger
	service *service.AnimeService
}

func NewAnimeMux(logger *zap.SugaredLogger, service *service.AnimeService) *AnimeMux {
	return &AnimeMux{
		logger:  logger,
		service: service,
	}
}

func (am *AnimeMux) AssignHandlers(routerGroup *gin.RouterGroup) {
	// /api/animes/{anime_id}?mark=[fav,unfav,clear]
	routerGroup.GET("/animes/:anime_id", am.UpdateMark)
	routerGroup.GET("/animes/:anime_id/", am.UpdateMark)

	// /api/animes?page=N
	routerGroup.GET("/animes", am.getPage)
}

func (am *AnimeMux) getPage(c *gin.Context) {
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil || page < 1 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	animes := am.service.GetPage(int(page))
	c.Status(http.StatusOK)
	json.NewEncoder(c.Writer).Encode(animes)
}

func (am *AnimeMux) UpdateMark(c *gin.Context) {
	am.logger.Infow("UpdateMark")
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
