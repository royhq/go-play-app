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
	var (
		errResponse   errorResponse
		validationErr *ValidationError
		cmdError      *CommandError
	)

	switch {
	case errors.As(e, &validationErr):
		errResponse = errorResponse{
			StatusCode: http.StatusBadRequest,
			Code:       validationErr.Code(),
			Msg:        validationErr.Msg,
		}

	case errors.As(e, &cmdError):
		errResponse = errorResponse{
			StatusCode: http.StatusInternalServerError,
			Code:       cmdError.Code,
			Msg:        cmdError.Msg,
		}

	default:
		errResponse = errorResponse{
			StatusCode: http.StatusInternalServerError,
			Code:       "internal_error",
			Msg:        "unexpected error",
		}
	}

	h.logger.ErrorContext(ctx, "error handled",
		slog.Any("error", e),
		slog.Any("response", errResponse),
		slog.Int("status_code", errResponse.StatusCode),
		slog.Any("stack", stack(e)),
	)

	resp.JSONResponse(w, errResponse.StatusCode, errResponse)
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
