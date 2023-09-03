package main

import (
	"anicomend/internal/modules/csv"
	"anicomend/internal/modules/http"
	"anicomend/internal/modules/logger"
	"anicomend/service"
	"net"
	"strconv"

	"go.uber.org/fx"
)

const (
	datasetPath = "anime_cleaned.csv"
	defaultHost = ""
	defaultPort = 8080
)

func main() {
	app := fx.New(
		fx.Provide(
			fx.Annotate(
				func() string { return datasetPath },
				fx.ResultTags(`name:"dataset-path"`),
			),
			service.NewAnimeService,
		),
		logger.Module,
		csv.Module,
		http.Module,
		fx.Invoke(func(s *http.Server) {
			s.Run(net.JoinHostPort(defaultHost, strconv.FormatInt(defaultPort, 10)))
		}),
	)

	app.Run()
}
