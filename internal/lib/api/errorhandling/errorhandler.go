package errorhandling

import (
	"blog/internal/constants/servererror"
	resp "blog/internal/lib/api/response"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

func HandleErrors(w http.ResponseWriter, r *http.Request, e error) bool {
	switch {
	case errors.Is(e, servererror.InvalidJson):
		render.JSON(w, r, resp.Error("failed to decode JSON", http.StatusBadRequest))
	case errors.Is(e, servererror.ResourceNotFound):
		render.JSON(w, r, resp.Error("resource not found", http.StatusNotFound))
	case errors.Is(e, servererror.InternalError):
		render.JSON(w, r, resp.Error("internal server error", http.StatusInternalServerError))
	case errors.Is(e, servererror.InvalidCrerdentials):
		render.JSON(w, r, resp.Error("invalid password or username", http.StatusUnauthorized))
	}
	return e != nil
}
