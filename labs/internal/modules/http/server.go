package http

import (
	"anicomend/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	logger       *zap.SugaredLogger
	engine       *gin.Engine
	animeService *service.AnimeService
}

func NewServer(logger *zap.SugaredLogger, apiMux *ApiMux, viewMux *ViewMux, animeService *service.AnimeService) *Server {
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

	apiMux.AssignHandlers(r.Group("/api"))
	viewMux.AssignHandlers(r.Group("/"))

	r.StaticFile("favicon.ico", "./static/favicon.ico")

	return &s
}

func (s *Server) Run(address string) error {
	s.logger.Infow("Run", "address", address)

	return s.engine.Run(address)
}
