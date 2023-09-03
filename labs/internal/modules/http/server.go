package http

import (
	"anicomend/internal/modules/http/dto"
	"anicomend/internal/modules/http/layout"
	"anicomend/service"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type Server struct {
	logger       *zap.SugaredLogger
	mux          *http.ServeMux
	animeService *service.AnimeService
}

func NewServer(logger *zap.SugaredLogger, am *AnimeMux, animeService *service.AnimeService) *Server {
	s := Server{
		logger:       logger,
		mux:          http.NewServeMux(),
		animeService: animeService,
	}

	// /api/animes?page=N
	s.mux.Handle("/api/animes/", http.StripPrefix("/api/animes", am))

	s.mux.Handle("/animes/", http.HandlerFunc(s.getAnimesPage))

	s.mux.Handle("/", http.RedirectHandler("/animes/", http.StatusMovedPermanently))

	return &s
}

func (s *Server) Run(address string) error {
	s.logger.Infow("Run", "address", address)

	return http.ListenAndServe(address, s.mux)
}

func (s *Server) getAnimesPage(w http.ResponseWriter, r *http.Request) {
	const maxPages = 9

	currentPage, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if err != nil {
		currentPage = 1
	}

	totalPages := s.animeService.GetTotalPages()
	animes := s.animeService.GetPage(int(currentPage))

	animeDTOs := make([]dto.AnimeDTO, len(animes))
	for i, anime := range animes {
		animeDTOs[i] = dto.NewAnimeDTO(anime)
	}

	layout.HomeLayout(w, layout.HomeLayoutParams{
		Animes:    animeDTOs,
		Pages:     layout.FormatPages(totalPages, maxPages, int(currentPage)),
		FirstPage: 1,
		LastPage:  totalPages,
	})
}
