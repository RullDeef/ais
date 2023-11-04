package layout

import (
	"anicomend/internal/modules/http/dto"
	"html/template"
	"io"
)

var recomendationsTemplate = template.Must(template.ParseFiles(
	"internal/modules/http/html/layout.tmpl",
	"internal/modules/http/html/filters.tmpl",
	"internal/modules/http/html/recomendations.tmpl",
	"internal/modules/http/html/anime_card.tmpl",
))

type RecomendationsLayoutParams struct {
	Animes        []dto.AnimeDTO
	Pages         []Page
	FirstPage     int
	LastPage      int
	NoPreferences bool
	FilterParams  FilterParams
}

func RecomendationsLayout(w io.Writer, params RecomendationsLayoutParams) error {
	return recomendationsTemplate.Execute(w, params)
}
