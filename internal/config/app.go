package config

import (
	"blog/config"
	"blog/pkg/db"
	"blog/pkg/logger"
	"database/sql"
	"log/slog"
)

type AppConfig struct {
	DB	   *sql.DB
	Logger *slog.Logger
	Config *config.Config
	User   *UserConfig
}

func NewAppConfig() *AppConfig {
	cfg := config.Cfg
	db_ := db.Connect(cfg.StoragePath)
	logger := logger.SetupLogger(cfg.Env)

	userConfig := NewUserConfig(db_, logger)

	return &AppConfig{
		DB:		db_,
		Logger: logger,
		Config: cfg,
		User:	userConfig,
	}
}
