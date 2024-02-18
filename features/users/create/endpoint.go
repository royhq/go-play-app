package create

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	resp "go-play-app/infra/http/response"
)

type (
	CommandHandlerFunc      func(context.Context, Command) (CommandOutput, error)
	CommandErrorHandlerFunc func(http.ResponseWriter, error)
)

type EndpointHandler struct {
	handleCreateUser   CommandHandlerFunc
	handleCommandError CommandErrorHandlerFunc
}

func (h *EndpointHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.BadRequest(w, "bad request")
		return
	}

	cmd := Command{
		Name: req.Name,
		Age:  req.Age,
	}

	out, err := h.handleCreateUser(ctx, cmd)
	if err != nil {
		h.handleCommandError(w, err)
		return
	}

	resp.Created(w, toResponse(out))
}

func NewEndpointHandler(cmdHandler CommandHandlerFunc, errHandler CommandErrorHandlerFunc) http.Handler {
	return &EndpointHandler{
		handleCreateUser:   cmdHandler,
		handleCommandError: errHandler,
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
