package posts

import (
	"blog/internal/lib/api/errorhandling"
	resp "blog/internal/lib/api/response"
	"blog/internal/persistence"
	"blog/internal/services"
	"blog/internal/types"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type GetAllResponse struct {
	resp.Response
	Posts []types.PostResponse
}

func GetAll(logger *slog.Logger,
	userRepo persistence.UserRepo,
	postRepo persistence.PostsRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.posts.GetAll"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		postService := services.NewPosts(logger, postRepo, userRepo)

		posts, err := postService.GetAll()
		if errorhandling.HandleErrors(w, r, err) {
			return
		}

		render.JSON(w, r, GetAllResponse{
			Response: resp.Ok(),
			Posts:    *posts,
		})
	}
}
