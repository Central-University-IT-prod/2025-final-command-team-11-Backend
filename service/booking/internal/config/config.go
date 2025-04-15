package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"REDACTED/team-11/backend/booking/pkg/postgres"
	"REDACTED/team-11/backend/booking/pkg/redis"
)

type Config struct {
	ServerPort      int    `env:"SERVER_PORT" env-default:"8080"`
	LogLevel        string `env:"LOG_LEVEL" env-default:"info"`
	JWTSecret       string `env:"JWT_SECRET" env-default:"root"`
	CoffeeIdBaseUrl string `env:"COFFEE_ID_BASE_URL" env-default:"http://localhost:8090"`
	PostgresConfig  postgres.Config
	RedisConfig     redis.Config
}

func Get() (Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
