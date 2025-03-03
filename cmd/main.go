package main

import (
	"blog/internal/deps"
	"blog/internal/middleware"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	app := deps.NewAppDeps()

	r := chi.NewRouter()

	r.Post("/register", app.User.Handler.Register)
	r.Post("/login", app.User.Handler.Login)
	r.Post("/refresh", app.User.Handler.Refresh)

	r.Group(func(pr chi.Router) {
		pr.Use(middleware.JwtMiddleware)
		pr.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("protected route"))
		})
		pr.Post("/avatar", app.User.Handler.UploadAvatar)
	})

	srv := http.Server{
		Addr:	 app.Config.Address,
		Handler: r,
	}

	app.Logger.Info("starting server on:", slog.String("address", app.Config.Address))

	srv.ListenAndServe()
}
