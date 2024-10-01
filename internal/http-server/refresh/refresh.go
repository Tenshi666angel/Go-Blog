package refresh

import (
	"blog/internal/constants/servererror"
	"blog/internal/lib/logger/sl"
	"blog/internal/persistence"
	"blog/internal/services"
	"blog/internal/token"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Status int
	Token  string `json:"token"`
}

func New(logger *slog.Logger, userRepo persistence.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.refresh.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		cookie, err := r.Cookie("tasty_cookies")
		if err != nil {
			logger.Error("cookie not found", sl.Err(err))
			panic(servererror.ResourceNotFound)
		}

		authService := services.NewAuth(logger, userRepo)

		tokenPair := authService.Refresh(cookie.Value)

		token.SetToCookie(tokenPair.RefreshToken, w)

		render.JSON(w, r, Response{
			Status: http.StatusOK,
			Token:  tokenPair.AccesToken,
		})
	}
}
