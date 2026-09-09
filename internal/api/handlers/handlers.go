package api_handlers

import (
	"context"
	"log/slog"

	"github.com/emp2ty0/coal-mine/internal/domain"
)

type HTTPHandlers struct {
	enterprise *domain.Enterprise
	logger     *slog.Logger
	ctx        context.Context
}

func NewHTTPHandlers(enterprise *domain.Enterprise, logger *slog.Logger, ctx context.Context) *HTTPHandlers {
	return &HTTPHandlers{
		enterprise: enterprise,
		logger:     logger,
		ctx:        ctx,
	}
}
