package dto

import (
	"anicomend/service"
	"fmt"
)

func ChatMessageToHTML(msg *service.ChatMessage) string {
	msgClass := "bot-message"
	if msg.FromUser {
		msgClass = "user-message"
	}
	html := fmt.Sprintf("<div class=\"mb-2 p-2 d-inline-block %s\">%s</div><br/>", msgClass, msg.Text)
	if msg.Animes != nil {
		html += "<div class=\"anime-row\">"
		for _, anime := range msg.Animes {
			html += fmt.Sprintf("<div>anime '%s'</div>", anime.Title)
		}
		html += "</div><br/>"
	}
	return html
}

func ChatErrorToHTML(err error) string {
	return fmt.Sprintf("<div class=\"mb-2 p-2 d-inline-block bot-error\">%s</div><br/>", err.Error())
}
