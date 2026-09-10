package api_handlers

import (
	"log/slog"

	"github.com/emp2ty0/coal-mine/internal/domain"
)

type HTTPHandlers struct {
	enterprise *domain.Enterprise
	logger     *slog.Logger
}

func NewHTTPHandlers(enterprise *domain.Enterprise, logger *slog.Logger) *HTTPHandlers {
	return &HTTPHandlers{
		enterprise: enterprise,
		logger:     logger,
	}
}
