package register

import (
	"blog/internal/lib/api/errorhandling"
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
			render.JSON(w, r, resp.Error("invalid json", http.StatusBadRequest))
		}

		authService := services.NewAuth(logger, userRepo)

		user, err := authService.Register(req)
		if errorhandling.HandleErrors(w, r, err) {
			return
		}

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			User:     *user,
		})
	}
}
