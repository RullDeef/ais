package nlp

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

var (
	ErrServiceUnavailable = errors.New("nlp service unavailable")
)

type NLPService struct {
	baseURL string
	logger  *zap.SugaredLogger
}

func New(baseURL string, logger *zap.SugaredLogger) *NLPService {
	return &NLPService{
		baseURL: baseURL,
		logger:  logger,
	}
}

func (s *NLPService) Request(message string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s?query=%s", s.baseURL, url.PathEscape(message)))
	if err != nil {
		s.logger.Errorw("failed to GET", "err", err, "message", message, "resp", resp)
		return "", ErrServiceUnavailable
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorw("failed to ReadAll", "err", err)
		return "", ErrServiceUnavailable
	}

	return string(bytes), nil
}
