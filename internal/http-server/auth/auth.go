package auth

import (
	"blog/internal/constants/servererror"
	"blog/internal/lib/api/errorhandling"
	"blog/internal/lib/logger/sl"
	"blog/internal/persistence"
	"blog/internal/services"
	"blog/internal/token"
	"blog/internal/types"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Status int
	Token  string `json:"token"`
}

func New(logger *slog.Logger, userRepo persistence.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.auth.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req types.User

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logger.Error("failed to decode json", sl.Err(err))
			panic(servererror.InvalidJson)
		}

		authService := services.NewAuth(logger, userRepo)

		tokenPair, err := authService.LogIn(req)
		if errorhandling.HandleErrors(w, r, err) {
			return
		}

		token.SetToCookie(tokenPair.RefreshToken, w)

		render.JSON(w, r, Response{
			Status: http.StatusOK,
			Token:  tokenPair.AccesToken,
		})
	}
}
