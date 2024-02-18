package create

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	resp "go-play-app/infra/http/response"
)

type EndpointErrorHandler struct {
	logger *slog.Logger
}

func (h *EndpointErrorHandler) HandleError(ctx context.Context, w http.ResponseWriter, e error) {
	h.logger.ErrorContext(ctx, "error handled", "error", e)

	if errors.Is(e, ErrValidation) {
		resp.BadRequest(w, e.Error())
		return
	}

	resp.InternalError(w, e.Error())
}

func NewEndpointErrorHandler(logger *slog.Logger) *EndpointErrorHandler {
	return &EndpointErrorHandler{logger: logger}
}
