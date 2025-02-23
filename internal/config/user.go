package config

import (
	"blog/internal/user"
	"database/sql"
	"log/slog"
)

type UserConfig struct {
	Handler *user.UserHandler
	Repo    user.UserRepo
	Service user.UserService 
}

func NewUserConfig(db *sql.DB, logger *slog.Logger) *UserConfig {
	repo := user.NewRepo(db, logger)
	service := user.NewService(logger, repo)
	handler := user.NewHandler(logger, service)	

	return &UserConfig{
		Handler: handler,
		Repo:    repo,
		Service: service,
	}
}
