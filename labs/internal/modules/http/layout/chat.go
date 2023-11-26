package layout

import (
	"html/template"
	"io"
)

var chatTemplate = template.Must(template.ParseFiles(
	"internal/modules/http/html/layout.tmpl",
	"internal/modules/http/html/chat_bot.tmpl",
	"internal/modules/http/html/anime_card.tmpl",
))

type ChatLayoutParams struct {
}

func ChatLayout(w io.Writer, params ChatLayoutParams) error {
	return chatTemplate.Execute(w, params)
}
