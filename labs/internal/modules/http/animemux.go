package http

import (
	"anicomend/service"
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type AnimeMux struct {
	http.ServeMux

	logger  *zap.SugaredLogger
	service *service.AnimeService
}

func NewAnimeMux(logger *zap.SugaredLogger, service *service.AnimeService) *AnimeMux {
	am := AnimeMux{
		ServeMux: *http.NewServeMux(),
		logger:   logger,
		service:  service,
	}

	// /api/animes?page=N
	am.Handle("/",
		panicWrapperMiddleware(
			loggerMiddleware(
				am.logger,
				http.HandlerFunc(am.getPage),
			),
		),
	)

	return &am
}

func (am *AnimeMux) getPage(w http.ResponseWriter, r *http.Request) {
	page_str := r.FormValue("page")
	page, err := strconv.ParseInt(page_str, 10, 64)
	if err != nil || page < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	animes := am.service.GetPage(int(page))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(animes)
}
