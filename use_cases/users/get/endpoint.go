package get

import (
	"context"
	"net/http"

	"github.com/royhq/go-play-app/shared/domain"
	"github.com/royhq/go-play-app/shared/infra/http/endpoints"
	resp "github.com/royhq/go-play-app/shared/infra/http/response"
)

type (
	QueryHandler interface {
		Handle(context.Context, Query) (QueryOutput, error)
	}

	ErrorHandler interface {
		HandleError(context.Context, http.ResponseWriter, error)
	}
)

type EndpointHandler struct {
	getUser      QueryHandler
	errorHandler ErrorHandler
}

func (h *EndpointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.PathValue("user_id")
	if userID == "" {
		resp.APIError(w, &endpoints.APIError{
			HTTPStatusCode: http.StatusBadRequest,
			Code:           "bad_request",
			Msg:            "user_id is required",
		})
		return
	}

	query := Query{UserID: domain.UserID(userID)}

	out, err := h.getUser.Handle(ctx, query)
	if err != nil {
		h.errorHandler.HandleError(ctx, w, err)
		return
	}

	resp.Ok(w, createResponse(out))
}

func createResponse(out QueryOutput) response {
	return response{
		UserID: string(out.User.ID),
		Name:   out.User.Name,
	}
}

type response struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}
