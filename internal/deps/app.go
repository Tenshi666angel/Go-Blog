package deps

import (
	"blog/config"
	"blog/internal/dbxutils"
	"blog/pkg/db"
	"blog/pkg/dbx"
	"blog/pkg/logger"
	"database/sql"
	"log/slog"
)

type AppDeps struct {
	DB	   *sql.DB
	Logger *slog.Logger
	Config *config.Config
	User   *UserDeps
}

func NewAppDeps() *AppDeps {
	cfg := config.Cfg
	db_ := db.Connect(cfg.StoragePath)
	logger := logger.SetupLogger(cfg.Env)
	dbxCommitter := dbxutils.NewDbxCommiter(*dbx.InitDbx(cfg.DbxToken))
	userConfig := NewUserDeps(db_, logger, *dbxCommitter)

	return &AppDeps{
		DB:		db_,
		Logger: logger,
		Config: cfg,
		User:	userConfig,
	}
}
