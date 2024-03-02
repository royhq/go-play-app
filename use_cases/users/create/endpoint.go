package create

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/royhq/go-play-app/shared/infra/http/endpoints"
	resp "github.com/royhq/go-play-app/shared/infra/http/response"
)

type (
	CommandHandlerFunc func(context.Context, Command) (CommandOutput, error)

	ErrorHandler interface {
		HandleError(context.Context, http.ResponseWriter, error)
	}
)

type EndpointHandler struct {
	createUser   CommandHandlerFunc
	errorHandler ErrorHandler
}

func (h *EndpointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.APIError(w, &endpoints.APIError{
			HTTPStatusCode: http.StatusBadRequest,
			Code:           "bad_request",
			Msg:            "bad request",
		})
		return
	}

	cmd := Command(req)

	out, err := h.createUser(ctx, cmd)
	if err != nil {
		h.errorHandler.HandleError(ctx, w, err)
		return
	}

	resp.Created(w, toResponse(out))
}

func NewEndpointHandler(cmdHandler CommandHandlerFunc, errHandler ErrorHandler) http.Handler {
	return &EndpointHandler{
		createUser:   cmdHandler,
		errorHandler: errHandler,
	}
}

type request struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type response struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

func toResponse(out CommandOutput) response {
	return response{
		ID:        string(out.CreatedUser.ID),
		Name:      out.CreatedUser.Name,
		Age:       out.CreatedUser.Age,
		CreatedAt: out.CreatedUser.Created,
	}
}
