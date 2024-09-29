package register

import (
	"blog/internal/constants/servererror"
	resp "blog/internal/lib/api/response"
	"blog/internal/lib/logger/sl"
	"blog/internal/persistence"
	"blog/internal/services"

	"blog/internal/types"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	User types.User `json:"user"`
}

func New(logger *slog.Logger, userRepo persistence.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.register.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req types.User

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logger.Error("failed to decode request body", sl.Err(err))
			panic(servererror.InvalidJson)
		}

		authService := services.NewAuth(logger, userRepo)

		user := authService.Register(req)

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			User:     *user,
		})
	}
}
