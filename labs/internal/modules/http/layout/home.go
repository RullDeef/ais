package layout

import (
	"anicomend/internal/modules/http/dto"
	"html/template"
	"io"
)

var homeTemplate = template.Must(template.ParseFiles(
	"internal/modules/http/html/layout.tmpl",
	"internal/modules/http/html/filters.tmpl",
	"internal/modules/http/html/home.tmpl",
	"internal/modules/http/html/anime_card.tmpl",
))

type HomeLayoutParams struct {
	Animes         []dto.AnimeDTO
	Pages          []Page
	FirstPage      int
	LastPage       int
	SearchQuery    string
	IsSearchResult bool
	FilterParams   FilterParams
}

func HomeLayout(w io.Writer, params HomeLayoutParams) error {
	return homeTemplate.Execute(w, params)
}
