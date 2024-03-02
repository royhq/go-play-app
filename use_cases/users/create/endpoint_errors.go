package create

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/royhq/go-play-app/shared/infra/http/endpoints"
	resp "github.com/royhq/go-play-app/shared/infra/http/response"
)

type EndpointErrorHandler struct {
	logger *slog.Logger
}

func (h *EndpointErrorHandler) HandleError(ctx context.Context, w http.ResponseWriter, e error) {
	var (
		errResponse   *endpoints.APIError
		validationErr *ValidationError
		cmdError      *CommandError
	)

	switch {
	case errors.As(e, &validationErr):
		errResponse = &endpoints.APIError{
			HTTPStatusCode: http.StatusBadRequest,
			Code:           validationErr.Code(),
			Msg:            validationErr.Msg,
		}

	case errors.As(e, &cmdError):
		errResponse = &endpoints.APIError{
			HTTPStatusCode: http.StatusInternalServerError,
			Code:           cmdError.Code,
			Msg:            cmdError.Msg,
		}

	default:
		errResponse = &endpoints.APIError{
			HTTPStatusCode: http.StatusInternalServerError,
			Code:           "internal_error",
			Msg:            "unexpected error",
		}
	}

	h.logger.ErrorContext(ctx, "error handled",
		slog.Any("error", e),
		slog.Any("response", errResponse),
		slog.Int("status_code", errResponse.HTTPStatusCode),
		slog.Any("stack", stack(e)),
	)

	resp.JSONResponse(w, errResponse.HTTPStatusCode, errResponse)
}

func NewEndpointErrorHandler(logger *slog.Logger) *EndpointErrorHandler {
	return &EndpointErrorHandler{logger: logger}
}

func stack(err error) []string {
	errs := make([]string, 0)

	for e := err; e != nil; e = errors.Unwrap(e) {
		errs = append(errs, e.Error())
	}

	return errs
}
