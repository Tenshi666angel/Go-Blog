package deps

import (
	"blog/config"
	"blog/pkg/db"
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

	userConfig := NewUserDeps(db_, logger)

	return &AppDeps{
		DB:		db_,
		Logger: logger,
		Config: cfg,
		User:	userConfig,
	}
}
