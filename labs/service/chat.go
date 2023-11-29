package service

import (
	"anicomend/model"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrServiceUnavailable = errors.New("nlp service unavailable")
)

type ChatService struct {
	history []*ChatMessage

	animeService *AnimeService
	logger       *zap.SugaredLogger
}

type ChatMessage struct {
	FromUser bool
	Text     string
	Animes   []*model.Anime
}

func NewChatService(logger *zap.SugaredLogger, animeService *AnimeService) *ChatService {
	return &ChatService{
		history:      nil,
		animeService: animeService,
		logger:       logger,
	}
}

func (cs *ChatService) PostMessage(message string, onDone func(*ChatMessage, error)) {
	cs.logger.Infow("PostMessage", "message", message)

	cs.history = append(cs.history, &ChatMessage{
		FromUser: true,
		Text:     message,
		Animes:   nil,
	})

	resp, err := cs.buildResponse(message)
	onDone(resp, err)
}

func (cs *ChatService) buildResponse(message string) (*ChatMessage, error) {
	return nil, ErrServiceUnavailable
}
