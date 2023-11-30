package service

import (
	"anicomend/internal/modules/nlp"
	"anicomend/model"
	"errors"
	"slices"

	"go.uber.org/zap"
)

var (
	ErrServiceUnavailable = errors.New("nlp service unavailable")
)

type ChatService struct {
	history []*ChatMessage

	animeService *AnimeService
	nlpService   *nlp.NLPService
	logger       *zap.SugaredLogger
}

type ChatMessage struct {
	FromUser bool
	Text     string
	Animes   []*model.Anime
}

func NewChatService(logger *zap.SugaredLogger, animeService *AnimeService, nlpService *nlp.NLPService) *ChatService {
	return &ChatService{
		history:      nil,
		animeService: animeService,
		nlpService:   nlpService,
		logger:       logger,
	}
}

func (cs *ChatService) GetHistory() []*ChatMessage {
	return slices.Clone(cs.history)
}

func (cs *ChatService) PostMessage(message string) (*ChatMessage, error) {
	cs.logger.Infow("PostMessage", "message", message)

	cs.history = append(cs.history, &ChatMessage{
		FromUser: true,
		Text:     message,
		Animes:   nil,
	})

	resp, err := cs.buildResponse(message)
	if err == nil {
		cs.history = append(cs.history, resp)
	}

	return resp, err
}

func (cs *ChatService) buildResponse(message string) (*ChatMessage, error) {
	resp, err := cs.nlpService.Request(message)
	if err != nil {
		cs.logger.Errorw("nlp request", "err", err, "message", message)
		return nil, ErrServiceUnavailable
	}

	// TODO: parse tags

	return &ChatMessage{
		FromUser: false,
		Text:     resp,
	}, nil
}
