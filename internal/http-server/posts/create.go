package posts

import (
	"blog/internal/lib/api/errorhandling"
	resp "blog/internal/lib/api/response"
	"blog/internal/lib/logger/sl"
	"blog/internal/persistence"
	"blog/internal/services"
	"blog/internal/types"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type CreateResponse struct {
	resp.Response
	Post types.PostResponse
}

func Create(logger *slog.Logger,
	postRepo persistence.PostsRepo,
	userRepo persistence.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.posts.Create"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		authHeader := r.Header.Get("Authorization")

		tokenString := strings.Split(authHeader, " ")[1]

		postService := services.NewPosts(logger, postRepo, userRepo)

		var req types.PostRequest

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logger.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid json", http.StatusBadRequest))
		}

		post, err := postService.Create(req, tokenString)
		if errorhandling.HandleErrors(w, r, err) {
			return
		}

		render.JSON(w, r, CreateResponse{
			Response: resp.Ok(),
			Post:     *post,
		})
	}
}
