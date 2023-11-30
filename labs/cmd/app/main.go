package main

import (
	"anicomend/internal/modules/csv"
	"anicomend/internal/modules/http"
	"anicomend/internal/modules/logger"
	"anicomend/internal/modules/nlp"
	"anicomend/service"
	"net"
	"os"
	"strconv"

	"go.uber.org/fx"
)

const (
	datasetPath = "anime_cleaned.csv"
	defaultHost = ""
	defaultPort = 8080
)

var (
	nlpBaseURL = getEnv("NLP_BASE_URL", "http://localhost:8085/")
)

func main() {
	app := fx.New(
		fx.Provide(
			fx.Annotate(
				func() string { return datasetPath },
				fx.ResultTags(`name:"dataset-path"`),
			),
			fx.Annotate(
				func() string { return nlpBaseURL },
				fx.ResultTags(`name:"nlp-base-url"`),
			),
		),
		logger.Module,
		csv.Module,
		http.Module,
		nlp.Module,
		service.Module,
		fx.Invoke(func(s *http.Server) {
			s.Run(net.JoinHostPort(defaultHost, strconv.FormatInt(defaultPort, 10)))
		}),
	)

	app.Run()
}

func getEnv(name, def string) string {
	val := os.Getenv(name)
	if val == "" {
		val = def
	}
	return val
}
