package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env			string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`
	Address		string `yaml:"address"`
	JwtKey      string `yaml:"jwt_key"`
	DbxToken    string `yaml:"dbx_token"`
}

func MustLoad() *Config {
	wd, err := os.Getwd()
    if err != nil {
        panic(err.Error())
    }

    configPath := filepath.Join(wd, "config", "config.yaml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exists", configPath)
		panic(err.Error())
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error with loading config: %s", err.Error())
		panic(err.Error())
	}

	return &cfg
}

var Cfg = MustLoad()
