package layout

import (
	"anicomend/internal/modules/http/dto"
	"html/template"
	"io"
)

var preferencesTemplate = template.Must(template.ParseFiles(
	"internal/modules/http/html/layout.tmpl",
	"internal/modules/http/html/preferences.tmpl",
	"internal/modules/http/html/anime_card.tmpl",
))

type PreferencesLayoutParams struct {
	Animes        []dto.AnimeDTO
	NoPreferences bool
}

func PreferencesLayout(w io.Writer, params PreferencesLayoutParams) error {
	return preferencesTemplate.Execute(w, params)
}
