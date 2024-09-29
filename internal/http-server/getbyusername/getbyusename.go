package getbyusername

import (
	"blog/internal/constants/servererror"
	resp "blog/internal/lib/api/response"
	"blog/internal/lib/logger/sl"
	"blog/internal/types"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Username string `json:"username"`
}

type Response struct {
	resp.Response
	User types.User `json:"user"`
}

type UserGetter interface {
	GetByUsername(username string) (*types.User, error)
}

func New(logger *slog.Logger, userGetter UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.getbyusername.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logger.Error("failed to decode json", sl.Err(err))
			panic(servererror.InvalidJson)
		}

		logger.Info("request body decoded", slog.Any("request", req))

		user, err := userGetter.GetByUsername(req.Username)
		if err != nil {
			logger.Error("failed to get user", sl.Err(err))
			panic(servererror.ResourceNotFound)
		}

		render.JSON(
			w, r, Response{
				Response: resp.Ok(),
				User:     *user,
			},
		)
	}
}
