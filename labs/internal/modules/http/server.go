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

type Server struct {
	logger       *zap.SugaredLogger
	engine       *gin.Engine
	animeService *service.AnimeService
}

func NewServer(logger *zap.SugaredLogger, animeMux *AnimeMux, animeService *service.AnimeService) *Server {
	r := gin.Default()
	s := Server{
		logger:       logger,
		engine:       r,
		animeService: animeService,
	}

	r.Use(func(ctx *gin.Context) {
		defer ctx.Header("Cache-Control", "no-cache")
		ctx.Next()
	})

	animeMux.AssignHandlers(r.Group("/api"))

	s.engine.GET("/animes", s.getAnimesPage)

	s.engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/animes")
	})

	return &s
}

func (s *Server) Run(address string) error {
	s.logger.Infow("Run", "address", address)

	return s.engine.Run(address)
}

func (s *Server) getAnimesPage(c *gin.Context) {
	const maxPages = 9

	currentPage, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		currentPage = 1
	}

	totalPages := s.animeService.GetTotalPages()
	animes := s.animeService.GetPage(int(currentPage))

	animeDTOs := make([]dto.AnimeDTO, len(animes))
	for i, anime := range animes {
		animeDTOs[i] = dto.NewAnimeDTO(anime, s.animeService.GetPreference(anime.Id))
	}

	c.Status(http.StatusOK)
	layout.HomeLayout(c.Writer, layout.HomeLayoutParams{
		Animes:    animeDTOs,
		Pages:     layout.FormatPages(totalPages, maxPages, int(currentPage)),
		FirstPage: 1,
		LastPage:  totalPages,
	})
}
