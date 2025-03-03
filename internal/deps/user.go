package deps

import (
	"blog/internal/dbxutils"
	"blog/internal/user"
	"database/sql"
	"log/slog"
)

type UserDeps struct {
	Handler     *user.UserHandler
	Repo        user.UserRepo
	Service     user.UserService
}

func NewUserDeps(db *sql.DB, logger *slog.Logger, dbxCommitter dbxutils.DbxCommiter) *UserDeps {
	repo := user.NewRepo(db, logger)
	service := user.NewService(logger, repo, dbxCommitter)
	handler := user.NewHandler(logger, service)	

	return &UserDeps{
		Handler:     handler,
		Repo:        repo,
		Service:     service,
	}
}
