package app

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log/slog"
)

type HTTPServerConfig struct {
	Address      string `env:"SERVICE_ADDRESS" env-default:"0.0.0.0:8080"` // Адрес и порт для HTTP сервера
	ReadTimeout  uint   `env:"SERVICE_READ_TIMEOUT" env-default:"30"`
	WriteTimeout uint   `env:"SERVICE_WRITE_TIMEOUT" env-default:"30"`
	IdleTimeout  uint   `env:"SERVICE_IDLE_TIMEOUT" env-default:"30"`
}

type Postgres struct {
	Conn        string `env:"POSTGRES_CONN"` // URL-строка для подключения к PostgreSQL
	AutoMigrate bool   `env:"POSTGRES_AUTO_MIGRATE" env-default:"false"`
	Migration   string `env:"POSTGRES_MIGRATION"`
}

type Config struct {
	HTTPServer HTTPServerConfig
	Postgres   Postgres
}

func mustLoadConfig(log slog.Logger) *Config {

	if err := godotenv.Load(); err != nil {
		log.Error("Error loading .env file")
	}

	config := new(Config)

	if err := cleanenv.ReadEnv(config); err != nil {
		log.Error("error reading config file: %s", err)
		return nil
	}

	return config
}
