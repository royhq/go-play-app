package create

import (
	"errors"
	"net/http"

	resp "go-play-app/infra/http/response"
)

func HandleCommandError(w http.ResponseWriter, e error) {
	if errors.Is(e, ErrValidation) {
		resp.BadRequest(w, e.Error())
		return
	}

	resp.InternalError(w, e.Error())
}
