package main

import (
	"blog/internal/config"
	"blog/internal/http-server/auth"
	"blog/internal/http-server/getbyusername"
	"blog/internal/http-server/middleware"
	"blog/internal/http-server/posts"
	"blog/internal/http-server/refresh"
	"blog/internal/http-server/register"
	"blog/internal/lib/logger/sl"
	"blog/internal/logger"
	"blog/internal/persistence/sqlite"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	conf := config.MustLoad()

	logger := logger.SetupLogger(conf.Env)

	storage, err := sqlite.New(conf.StoragePath)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
	}

	r := chi.NewRouter()

	r.Group(func(pr chi.Router) {
		pr.Use(middleware.JwtMiddleware)
		pr.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("protected route"))
		})
		pr.Post("/posts/create", posts.Create(logger, storage, storage))
		pr.Get("/posts/getall", posts.GetAll(logger, storage, storage))
		pr.Post("/posts/like", posts.Like(logger, storage, storage))
		pr.Post("/posts/unlike", posts.UnLike(logger, storage, storage))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	r.Post("/register", register.New(logger, storage))
	r.Post("/refresh", refresh.New(logger, storage))
	r.Post("/login", auth.New(logger, storage))
	r.Get("/getuser", getbyusername.New(logger, storage))

	logger.Info(fmt.Sprintf("Starting server on %s", conf.Address))

	srv := &http.Server{
		Addr:         conf.Address,
		Handler:      r,
		ReadTimeout:  conf.Timeout,
		WriteTimeout: conf.Timeout,
		IdleTimeout:  conf.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server", sl.Err(err))
	}

	logger.Error("Server stopped")
}
